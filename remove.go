package main

import (
	"fmt"
	"os"
	"path"
)

func removeImport(dir string, importPath string) (err error) {
	var mainRoot string
	mainRoot, err = gitGetRootDir(dir)
	if err != nil {
		return
	}

	targetPath := path.Join("vendor", importPath)
	if _, err = os.Stat(path.Join(mainRoot, targetPath)); err != nil {
		return
	}
	err = nil

	var importRoot string
	absoluteImportPath := path.Join(mainRoot, targetPath)
	importRoot, err = gitGetRootDir(absoluteImportPath)

	if err != nil {
		return
	}

	if mainRoot == importRoot {
		err = fmt.Errorf("%q is not a submodule of %q", absoluteImportPath, mainRoot)
		return
	}
	importRoot = importRoot[len(mainRoot)+1:]

	if importRoot != targetPath {
		fmt.Printf("Removing %q implies removing the whole %q submodule. Do you wish to continue ? (Y/n) : ", targetPath, importRoot)
		var response string
		if _, err = fmt.Scanln(&response); err != nil {
			if err.Error() == "unexpected newline" {
				err = nil
			} else {
				return
			}
		}
		if response == "" || response[:1] == "y" || response[:1] == "Y" {
			targetPath = importRoot
		} else {
			fmt.Println("Cancelled.")
			return
		}
	}

	//var output []byte
	_, err = gitRemoveSubmodule(mainRoot, targetPath)
	//fmt.Println(string(output))

	return
}
