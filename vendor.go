package main

import (
	"os"
	"path"

	"fmt"
	"strings"

	"runtime"
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

	var importRoot string
	importRoot, err = gitGetRootDir(pkgDir)
	if err != nil {
		return
	}
	importPath = importRoot[len(gopath)+5:] // removes "$GOPATH/src/"

	var commitHash, remoteURL string

	commitHash, err = gitGetCurrentCommitHash(importRoot)
	if err != nil {
		return
	}

	remoteURL, err = gitGetRemoteURI(importRoot, false)
	if err != nil {
		return
	}

	var mainRoot string
	mainRoot, err = gitGetRootDir(dir)
	if err != nil {
		return
	}

	targetPath := path.Join("vendor", importPath)
	if _, err = os.Stat(targetPath); err == nil {
		//err = fmt.Errorf("%q already exists", targetPath)
		// Already exists, skipping
		return
	}

	_, err = gitAddSubmodule(mainRoot, remoteURL, targetPath)
	if err != nil {
		// @TODO : proper logging with verbose option
		//fmt.Println(string(output))
		return
	}

	_, err = gitCheckoutCommit(path.Join(mainRoot, targetPath), commitHash)
	if err != nil {
		//fmt.Println(string(output))
		return
	}

	vendored = true
	return
}
