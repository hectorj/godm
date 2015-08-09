package main

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/tools/go/vcs"
)

func removeImport(dir string, importPath string) (err error) {
	gopath := os.Getenv("GOPATH")
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
	if _, err = os.Stat(targetPath); err != nil {
		return
	}
	err = nil

	switch mainVCS.Name {
	case "Git":
		_, err = gitRemoveSubmodule(absoluteMainRoot, targetPath)
	default:
		err = fmt.Errorf("Unsupported VCS for dir %q : %q", dir, mainVCS.Name)
	}

	return
}
