package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

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
	paths := []string(c.Args())
	if len(paths) == 0 {
		cli.ShowCommandHelp(c, "vendor")
		fatalErrorf("No paths given")
	}

	currentDir := getAbsoluteCurrentDirOrExit()

	excludePatterns := c.StringSlice("exclude")
	excludes := make([]*regexp.Regexp, len(excludePatterns))
	var err error
	for index, pattern := range excludePatterns {
		excludes[index], err = regexp.Compile(pattern)
		if err != nil {
			fatalErrorf("Invalid regexp pattern %q : %s", pattern, err.Error())
		}
	}

	files, err := listFiles(paths, c.Bool("r"), excludes)
	if err != nil {
		fatalErrorf("Error listing files to scan for imports : %s", err.Error())
	}

	imports, err := extractImports(files)
	if err != nil {
		fatalErrorf("Error scanning imports from files : %s", err.Error())
	}

	for _, importPath := range imports {

		err, ok := vendorImport(currentDir, importPath)
		if err != nil {
			Log.Error("%s : Failed (%s)", importPath, err.Error())
		} else if !ok {
			Log.Info("%s : Skipped", importPath)
		} else {
			Log.Notice("%s : OK", importPath)
			fmt.Println(importPath)
		}
	}
}

func listFiles(fileNames []string, recursive bool, excludes []*regexp.Regexp) ([]string, error) {
	files := make([]string, 0, len(fileNames))
	var err error

	for _, fileName := range fileNames {
		err = listFile(fileName, &files, recursive, "", excludes)
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func listFile(filename string, list *[]string, recursiveScanning bool, currentPath string, excludes []*regexp.Regexp) error {
	for _, excludeRegexp := range excludes {
		if excludeRegexp.MatchString(filename) {
			return nil
		}
	}

	fullPath := path.Join(currentPath, filename)
	Log.Debug(fullPath)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		err = fmt.Errorf("Error with file %q : %q\n", filename, err.Error())
		return err
	}
	if fileInfo.IsDir() {
		if recursiveScanning || currentPath == "" {
			fileInfos, err := ioutil.ReadDir(fullPath)
			if err != nil {
				return err
			}
			for _, fileInfo = range fileInfos {
				if strings.HasSuffix(fileInfo.Name(), ".go") || fileInfo.IsDir() {
					err = listFile(fileInfo.Name(), list, recursiveScanning, fullPath, excludes)
					if err != nil {
						return err
					}
				}
			}
		}
	} else {
		*list = append(*list, fullPath)
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
