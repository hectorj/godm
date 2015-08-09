package main

import (
	"os"
	"path"

	"fmt"
	"strings"

	"runtime"

	"golang.org/x/tools/go/vcs"
)

func vendorImport(dir string, importPath string) (err error, vendored bool) {
	gopath := os.Getenv("GOPATH")
	goroot := runtime.GOROOT()
	pkgDir := path.Join(gopath, "src", importPath)
	if _, err = os.Stat(pkgDir); err != nil {
		if os.IsNotExist(err) {
			pkgDir = path.Join(goroot, "src", importPath)
			if _, err = os.Stat(pkgDir); err != nil {
				if os.IsNotExist(err) {
					err = fmt.Errorf("Import path %q not found in GOPATH %q nor in GOROOT %q", importPath, gopath, goroot)
				} else {
					// It's a standard library package, skipping
					err = nil
				}
			}
		}
		return
	}

	if strings.HasPrefix(pkgDir, dir) {
		// It is a sub-package of the current one, skipping
		return
	}

	var (
		importVCS  *vcs.Cmd
		importRoot string
	)
	importVCS, importRoot, err = vcs.FromDir(pkgDir, gopath)
	if err != nil {
		return
	}
	importPath = importRoot[4:] // removes "src/"
	absoluteImportRoot := path.Join(gopath, importRoot)

	var commitHash, remoteURL string
	//var output []byte
	switch importVCS.Name {
	case "Git":
		commitHash, err = gitGetCurrentCommitHash(absoluteImportRoot)
		if err != nil {
			return
		}

		remoteURL, err = gitGetRemoteURI(absoluteImportRoot, false)
		if err != nil {
			return
		}
	default:
		// @TODO : mercurial subrepositories (if `go get` supports them)
		err = fmt.Errorf("Unsupported VCS for dir %q : %q", pkgDir, importVCS.Name)
		return
	}

	var (
		mainVCS  *vcs.Cmd
		mainRoot string
	)
	mainVCS, mainRoot, err = vcs.FromDir(dir, gopath)
	if err != nil {
		return
	}
	absoluteMainRoot := path.Join(gopath, mainRoot)

	targetPath := path.Join("vendor", importPath)
	if _, err = os.Stat(targetPath); err == nil {
		//err = fmt.Errorf("%q already exists", targetPath)
		// Already exists, skipping
		return
	}
	err = nil

	switch mainVCS.Name {
	case "Git":
		_, err = gitAddSubmodule(absoluteMainRoot, remoteURL, targetPath)
		if err != nil {
			// @TODO : proper logging with verbose option
			//fmt.Println(string(output))
			return
		}

		_, err = gitCheckoutCommit(path.Join(absoluteMainRoot, targetPath), commitHash)
		if err != nil {
			//fmt.Println(string(output))
			return
		}

	default:
		err = fmt.Errorf("Unsupported VCS for dir %q : %q", dir, mainVCS.Name)
		return
	}
	vendored = true
	return
}
