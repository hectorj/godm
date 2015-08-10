package gpm

import "errors"

var ErrOrphan = errors.New("Vendor does not have a parent project")

// Vendor represents a vendored project
type Vendor interface {
	LocalProject
	// GetParent returns the project vendoring this one.
	GetParent() Project
	// SetParent saves the reference to the project vendoring this one.
	SetParent(parent Project)
	// GetBaseImportPath returns the base import path (the one containing all the eventual subpackages)
	GetImportPath() string
	// GetProject returns the vendor's project
	GetProject() LocalProject
}

type vendor struct {
	LocalProject
	parent     Project
	importPath string
}

func (self *vendor) GetProject() LocalProject {
	return self.LocalProject
}

func (self *vendor) GetImportPath() string {
	return self.importPath
}

func (self *vendor) GetParent() Project {
	return self.parent
}

func (self *vendor) SetParent(parent Project) {
	self.parent = parent
}

type goGettableVendor struct {
	Vendor
}

//func (self *goGettableVendor) Copy() error {
//	var parent Project
//	if parent = self.GetParent(); parent == nil {
//		// Vendor cannot copy itself into its parent if it does not have one.
//		return ErrOrphan
//	}
//	importPath := self.GetImportPath()
//
//	//// `Go get` the package in a temporary GOPATH
//	gopath, err := ioutil.TempDir("", "gopath-gpm-"+importPath)
//	if err != nil {
//		return err
//	}
//	defer os.RemoveAll(gopath)
//	goBin, err := exec.LookPath("go")
//	if err != nil {
//		return err
//	}
//	cmd := exec.Command(goBin, "get", importPath)
//	cmd.Env = []string{`GOPATH="` + gopath + `"`, "GO15VENDOREXPERIMENT=1"}
//	if err = cmd.Run(); err != nil {
//		return err
//	}
//	////
//
//	targetPath := path.Join(parent.GetBaseDir(), "vendor", importPath)
//	if err = os.MkdirAll(targetPath, os.ModeDir); err != nil {
//		return err
//	}
//	if err = CopyDir(path.Join(gopath, "src", importPath), targetPath); err != nil {
//		return err
//	}
//	return nil
//}
