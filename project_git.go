package gpm

import (
	"path"

	"fmt"

	"github.com/hectorj/gpm/git"
)

type RemoteGitProject interface {
	Project
	// GetURI returns a fetchable URI for the repository
	GetGitURI() string
}

type LocalGitProject interface {
	LocalProject
	// GetReference returns the HEAD reference.
	// Can be a commit hash, a branch, a tag etc.
	// Must be usable with the `git checkout` command
	GetReference() (string, error)
	// GetRemote returns the RemoteGitProject if possible
	// Returns nil if there is no remote.
	GetRemote() (RemoteGitProject, error)
}

type remoteGitProject struct {
	uri string
}

func (self *remoteGitProject) GetGitURI() string {
	return self.uri
}

func (self *remoteGitProject) Install(destination string) (LocalProject, error) {
	destination = path.Clean(destination)
	if err := git.Clone(destination, self.GetGitURI()); err != nil {
		return nil, err
	}
	return NewGitProjectFromPath(destination)
}

type localGitProject struct {
	// Falls back as a project without VCL when git features are not usable
	ProjectNoVCL
	remote        RemoteGitProject
	remoteChecked bool
	reference     string
}

var _ LocalGitProject = (*localGitProject)(nil)

func NewGitProjectFromPath(path string) (*localGitProject, error) {
	gitBaseDir, err := git.GetRootDir(path)
	if err != nil {
		return nil, err
	}
	project := &localGitProject{
		ProjectNoVCL: *(NewProjectNoVCL(gitBaseDir)),
	}
	project.Recursive = true
	return project, nil
}

func NewlocalGitProjectFromURI(uri, reference string) *localGitProject {
	return &localGitProject{}
}

func (self *localGitProject) GetReference() (reference string, err error) {
	if self.reference == "" {
		self.reference, err = self.getReference()
	}
	return self.reference, err
}

func (self *localGitProject) getReference() (reference string, err error) {
	if self.GetBaseDir() == "" {
		return "master", nil
	}
	return git.GetCurrentCommitHash(self.GetBaseDir())
}

func (self *localGitProject) GetRemote() (RemoteGitProject, error) {
	if !self.remoteChecked {
		URI, err := git.GetRemoteURI(self.GetBaseDir())
		if err != nil {
			if err == git.ErrNoRemote {
				self.remoteChecked = true
				self.remote = nil
			} else {
				return nil, err
			}
		}
		self.remoteChecked = true
		self.remote = &remoteGitProject{
			uri: URI,
		}
	}
	return self.remote, nil

}

// AddVendor as a git submodule if possible, or else by simply copying it
func (self *localGitProject) AddVendor(importPath string, project Project) (Vendor, error) {
	vendors, err := self.GetVendors()
	if err != nil {
		return nil, err
	}
	if _, exists := vendors[importPath]; exists {
		return nil, ErrDuplicateVendor
	}

	v := &vendor{
		parent:     self,
		importPath: importPath,
	}
	relativeTargetPath := path.Join("vendor", importPath)
	absoluteTargetPath := path.Join(self.GetBaseDir(), relativeTargetPath)
	switch typedProject := project.(type) {
	case RemoteGitProject:
		git.AddSubmodule(self.GetBaseDir(), typedProject.GetGitURI(), relativeTargetPath)
		v.LocalProject, err = NewGitProjectFromPath(absoluteTargetPath)
		return v, err
	case LocalGitProject:
		remote, err := typedProject.GetRemote()
		if err != nil {
			return nil, err
		}
		if remote != nil {
			reference, err := typedProject.GetReference()
			if err != nil {
				return nil, err
			}

			err = git.AddSubmodule(self.GetBaseDir(), remote.GetGitURI(), relativeTargetPath)
			if err != nil {
				return nil, err
			}

			errorHandler := func() {
				git.RemoveSubmodule(self.GetBaseDir(), relativeTargetPath)
			}
			defer func() {
				if panicErr := recover(); panicErr != nil {
					errorHandler()
					panic(panicErr)
				}
			}()

			err = git.CheckoutCommit(absoluteTargetPath, reference)
			if err != nil {
				errorHandler()
				return nil, err
			}

			v.LocalProject, err = NewGitProjectFromPath(absoluteTargetPath)

			return v, err
		}
	}

	if v.LocalProject, err = project.Install(absoluteTargetPath); err != nil {
		return nil, err
	}

	return v, nil
}

// RemoveVendor by removing it as a git submodule if necessary, then deleting it from the file system

func (self *localGitProject) RemoveVendor(importPath string) error {
	vendors, err := self.GetVendors()
	if err != nil {
		return err
	}
	vendor, exists := vendors[importPath]
	if !exists {
		return ErrUnknownVendor
	}
	v := vendor.GetProject().(LocalGitProject)
	fmt.Println(v.GetBaseDir())
	switch vendor.GetProject().(type) {
	case LocalGitProject:
		fmt.Println("test1")
		targetPath := path.Join("vendor", importPath)
		err := git.RemoveSubmodule(self.GetBaseDir(), targetPath)
		if err != nil {
			return err
		}
		vendor.SetParent(nil)
		delete(self.Vendors, importPath)
		return nil
	default:
		fmt.Println("test3")
		return self.ProjectNoVCL.RemoveVendor(importPath)
	}
}
