package godm

import (
	"errors"
	"os"
	pathlib "path"
	"regexp"

	"runtime"

	"strings"

	"github.com/hectorj/godm/git"
)

var (
	ErrDuplicateVendor = errors.New("Vendor already added to this project")
	ErrUnknownVendor   = errors.New("Vendor doesn't exist in this project")
)

var DefaultExcludesRegexp = []*regexp.Regexp{regexp.MustCompile(`/?vendor/?`)}

// Project represents a Go project
type Project interface {
	// Install downloads/copies/generates the project at specified destination on the filesystem
	Install(destination string) (LocalProject, error)
}

// LocalProject represents a Project we have available on the local file system
type LocalProject interface {
	Project
	// GetVendors returns all the project's vendors as a map[importPath]Vendor
	GetVendors() (vendors map[string]Vendor, err error)
	// AddVendor adds a vendor to the project
	AddVendor(importPath string, project Project) (Vendor, error)
	// RemoveVendor removes a vendor from the project by import path
	RemoveVendor(importPath string) error
	// GetImports returns all the project's imports (as a set of import paths)
	GetImports() (importPaths Set, err error)
	// GetSubpackages returns all the project's subpackages, (excluding vendors).
	GetSubpackages() (subpackages Set, err error)
	// GetBaseDir returns the path to the base directory of the project
	GetBaseDir() string
}

var ErrImportPathNotFound = errors.New("Import path not found in GOPATH nor as Go-Gettable")
var ErrNotImplemented = errors.New("Feature not implemented yet")
var ErrStandardLibrary = errors.New("Import path is part of the standard library")

func NewProjectFromImportPath(importPath string) (Project, string, error) {
	// Checking if the import path is in the GOPATH
	// @TODO : support multiple GOPATHs
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")

goPathLoop:
	for _, gopath := range gopaths {
		if gopath != "" {
			fullpath := pathlib.Join(gopath, "src", importPath)

			info, err := os.Stat(fullpath)
			if err != nil {
				if os.IsNotExist(err) {
					continue goPathLoop
				} else {
					return nil, importPath, err
				}
			}
			if info.IsDir() {
				// The dir does appear in the GOPATH
				project, err := NewLocalProject(fullpath)
				canonicalImportPath := pathlib.Clean(strings.TrimSuffix(importPath, strings.TrimPrefix(fullpath, project.GetBaseDir())))
				return project, canonicalImportPath, err
			}
		}
	}

	// Checking if the import path is in the GOROOT
	fullpath := pathlib.Join(runtime.GOROOT(), "src", importPath)
	if info, err := os.Stat(fullpath); err == nil && info.IsDir() {
		// The dir does appear in the GOROOT
		return nil, importPath, ErrStandardLibrary
	}

	//	if result, _, _, err := IsGoGettable(importPath); err != nil {
	//		return nil, err
	//	} else {
	//		switch result {
	//		case NotGoGettable:
	//			return nil, ErrImportPathNotFound
	//		default:
	return nil, importPath, ErrNotImplemented //@TODO
	//		}
	//	}
}

// NewLocalProject takes a path and builds a Project object from it, determining if it is a Git/Mercurial/No VCL/etc. project
func NewLocalProject(path string) (LocalProject, error) {
	gitProject, err := NewGitProjectFromPath(path)
	if err == git.ErrNotAGitRepository {
		return NewProjectNoVCL(path), nil
	} else if err != nil {
		return nil, err
	}
	return gitProject, nil
}
