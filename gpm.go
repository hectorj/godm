package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type actionType int

const (
	actionVendor actionType = iota
	actionRemove
)

type app struct {
	currentDir       string
	action           actionType
	vendorParameters struct {
		recursiveScanning bool
		files             []string
	}
	removeParameters struct {
		importPath string
	}
}

func (self *app) parseCommandLineArguments() (err error) {
	flag.Parse()
	self.currentDir, err = filepath.Abs(path.Dir(os.Args[0]))
	if err != nil {
		return
	}

	actionName := flag.Arg(0)
	switch actionName {
	case "vendor":
		self.action = actionVendor

		err = self.parseVendorCommandLineArguments()
	case "remove":
		self.action = actionRemove

		err = self.parseRemoveCommandLineArguments()
	case "":
		err = errors.New("Missing action's argument")
	default:
		err = fmt.Errorf("Unknown action %q", actionName)
	}

	return
}

func (self *app) parseRemoveCommandLineArguments() error {
	self.removeParameters.importPath = flag.Arg(1)
	if self.removeParameters.importPath == "" {
		return errors.New("No import path given")
	}
	return nil
}

func (self *app) parseVendorCommandLineArguments() error {
	flag.BoolVar(&self.vendorParameters.recursiveScanning, "r", false, "Scan dirs recursively")
	flag.Parse()

	// take all args without the flags and the action name
	pathArgs := flag.Args()[1:]

	pathsCount := len(pathArgs)
	if pathsCount == 0 {
		return errors.New("No paths given")
	}

	// pre-allocation
	self.vendorParameters.files = make([]string, 0, pathsCount)
	var err error

	for _, fileName := range pathArgs {
		err = listFile(fileName, &self.vendorParameters.files, self.vendorParameters.recursiveScanning, true)
		if err != nil {
			return err
		}
	}
	return nil
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
				err = listFile(fileInfo.Name(), list, recursiveScanning, false)
				if err != nil {
					return err
				}
			}
		}
	} else {
		*list = append(*list, filename)
	}
	return nil
}
