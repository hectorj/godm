package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"path/filepath"

	"path"

	"github.com/codegangsta/cli"
)

func remove(c *cli.Context) {
	importPath := c.Args().First()
	if importPath == "" {
		cli.ShowCommandHelp(c, "remove")
		fatalErrorf("No import path given")
	}

	err := removeImport(getAbsoluteCurrentDirOrExit(), importPath, c.Bool("y"))
	if err != nil {
		fatalErrorf("Error removing import : %s", err.Error())
	}
}

func vendor(c *cli.Context) {
	paths := c.Args().Tail()
	if len(paths) == 0 {
		cli.ShowCommandHelp(c, "vendor")
		fatalErrorf("No paths given")
	}

	currentDir := getAbsoluteCurrentDirOrExit()

	files, err := listFiles(paths, c.Bool("r"))
	if err != nil {
		fatalErrorf("Error listing files to scan for imports : %s", err.Error())
	}

	imports, err := extractImports(files)
	if err != nil {
		fatalErrorf("Error scanning imports from files : %s", err.Error())
	}

	for _, importPath := range imports {

		fmt.Print(importPath, " : ")
		err, ok := vendorImport(currentDir, importPath)
		if err != nil {
			fmt.Print("Failed (", err.Error(), ")")
		} else if !ok {
			fmt.Print("Skipped")
		} else {
			fmt.Print("OK")
		}
		fmt.Print("\n")
	}
}

func listFiles(fileNames []string, recursive bool) ([]string, error) {
	files := make([]string, 0, len(fileNames))
	var err error

	for _, fileName := range fileNames {
		err = listFile(fileName, &files, recursive, true)
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func listFile(filename string, list *[]string, recursiveScanning, firstLevel bool) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		err = fmt.Errorf("Error with file %q : %q\n", filename, err.Error())
		return err
	}
	if fileInfo.IsDir() {
		if recursiveScanning || firstLevel {
			fileInfos, err := ioutil.ReadDir(filename)
			if err != nil {
				return err
			}
			for _, fileInfo = range fileInfos {
				if strings.HasSuffix(fileInfo.Name(), ".go") {
					err = listFile(fileInfo.Name(), list, recursiveScanning, false)
					if err != nil {
						return err
					}
				}
			}
		}
	} else {
		*list = append(*list, filename)
	}
	return nil
}

func getAbsoluteCurrentDirOrExit() string {
	currentDir, err := filepath.Abs(path.Dir(os.Args[0]))
	if err != nil {
		fatalErrorf("Error getting current absolute path : %s", err.Error())
	}
	return currentDir
}
