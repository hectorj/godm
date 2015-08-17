package godm

import (
	"os"
	"path"
	"regexp"
	"strings"
)

// ProjectNoVCL is the most baisc local project : without any Version Control System
type ProjectNoVCL struct {
	BaseDir        string
	Recursive      bool
	Subpackages    Set
	Imports        Set
	Vendors        map[string]Vendor
	GoFiles        Set
	PathExcludes   []*regexp.Regexp
	vendorsChecked bool
}

var _ Project = (*ProjectNoVCL)(nil)

func NewProjectNoVCL(path string) *ProjectNoVCL {
	return &ProjectNoVCL{
		BaseDir:      path,
		Recursive:    false,
		PathExcludes: DefaultExcludesRegexp, //@TODO : parametrize
	}
}

func (self *ProjectNoVCL) GetBaseDir() string {
	return self.BaseDir
}

func (self *ProjectNoVCL) GetVendors() (vendors map[string]Vendor, err error) {
	if !self.vendorsChecked {
		self.Vendors, err = self.getVendors()
		self.vendorsChecked = err == nil
	}
	return self.Vendors, err
}

func (self *ProjectNoVCL) getVendors() (map[string]Vendor, error) {
	vendorPath := path.Join(self.GetBaseDir(), "vendor")
	vendorGoFiles, err := listGoFiles(vendorPath, true, nil, true)
	if err != nil {
		return nil, err
	}
	if len(vendorGoFiles) == 0 {
		return nil, nil
	}
	vendorDirs := NewSet()
	for filePath := range vendorGoFiles {
		vendorDirs.Add(path.Dir(filePath))
	}
	result := make(map[string]Vendor, len(vendorDirs))
	for vendorDir := range vendorDirs {
		// @TODO : optimize. We can know if a project includes its subpackages, avoiding us processing them for nothing
		vendorProject, err := NewLocalProject(vendorDir)
		if err != nil {
			return nil, err
		}
		nonCanonicalPathPart := strings.TrimPrefix(vendorDir, vendorProject.GetBaseDir())
		importPath := strings.Trim(strings.TrimSuffix(strings.TrimPrefix(vendorProject.GetBaseDir(), vendorPath), nonCanonicalPathPart), "/")
		//fmt.Printf("\n\n%s\n%s\n%s\n\n", nonCanonicalPathPart, importPath, vendorProject.GetBaseDir())
		//fmt.Printf("%s ; %s ; %s\n\n", vendorPath, vendorProject.GetBaseDir(), importPath)
		if _, exists := result[importPath]; !exists {
			result[importPath] = &vendor{
				parent:       self,
				LocalProject: vendorProject,
				importPath:   importPath,
			}
		}
	}
	return result, nil
}

func (self *ProjectNoVCL) GetImports() (importPaths Set, err error) {
	if self.Imports == nil {
		self.Imports, err = self.getImports()
	}
	return self.Imports, err
}

func (self *ProjectNoVCL) getImports() (importPaths Set, err error) {
	var goFiles Set
	if goFiles, err = self.getGoFiles(); err != nil {
		return
	}
	return extractImports(goFiles)
}

func (self *ProjectNoVCL) GetSubpackages() (subpackages Set, err error) {
	if self.Subpackages == nil {
		self.Subpackages, err = self.getSubpackages()
	}
	return self.Subpackages, err
}

func (self *ProjectNoVCL) getSubpackages() (subpackages Set, err error) {
	var goFiles Set
	if goFiles, err = self.getGoFiles(); err != nil {
		return
	}
	subpackages = NewSet()
	for file := range goFiles {
		subpackages.Add(path.Dir(file))
	}
	return
}

func (self *ProjectNoVCL) getGoFiles() (gofiles Set, err error) {
	if self.GoFiles == nil {
		self.GoFiles, err = listGoFiles(self.GetBaseDir(), self.Recursive, self.PathExcludes, true)
	}
	return self.GoFiles, err
}

// AddVendor by simply copying it
func (self *ProjectNoVCL) AddVendor(importPath string, project Project) (Vendor, error) {
	vendors, err := self.GetVendors()
	if err != nil {
		return nil, err
	}
	if _, exists := vendors[importPath]; exists {
		return nil, ErrDuplicateVendor
	}

	targetPath := path.Join(self.GetBaseDir(), "vendor", importPath)
	localProject, err := project.Install(targetPath)
	if err != nil {
		return nil, err
	}

	v := &vendor{
		LocalProject: localProject,
		parent:       self,
		importPath:   importPath,
	}
	self.Vendors[importPath] = v
	return v, nil
}

// RemoveVendor by deleting it from the file system
func (self *ProjectNoVCL) RemoveVendor(importPath string) error {
	vendors, err := self.GetVendors()
	if err != nil {
		return err
	}
	vendor, exists := vendors[importPath]
	if !exists {
		return ErrUnknownVendor
	}

	targetPath := path.Join(self.GetBaseDir(), "vendor", vendor.GetImportPath())
	if err := os.RemoveAll(targetPath); err != nil {
		return err
	}
	vendor.SetParent(nil)
	delete(self.Vendors, importPath)
	RemoveSubdirsWithNoFiles(path.Join(self.GetBaseDir(), "vendor", strings.Split(vendor.GetImportPath(), "/")[0]))
	return nil
}

func (self *ProjectNoVCL) Install(destination string) (LocalProject, error) {
	destination = path.Clean(destination)
	if err := CopyDir(self.GetBaseDir(), destination); err != nil {
		return nil, err
	}
	return NewProjectNoVCL(destination), nil
}
