package main

import (
	"os"
	"path"

	"fmt"

	"github.com/codegangsta/cli"
	"github.com/hectorj/godm"
)

func clean(c *cli.Context) {
	if len(c.Args()) > 0 {
		cli.ShowCommandHelp(c, "clean")
		fatalErrorf("Too much arguments")
	}

	project, err := godm.NewLocalProject(path.Dir(os.Args[0]))
	if err != nil {
		fatalErrorf("Error building the current project : %s", err.Error())
	}
	Log.Debug("Project's Type : %T", project, project.GetBaseDir())
	Log.Debug("Project's Base Dir : %s", project.GetBaseDir())

	imports, err := project.GetImports()
	if err != nil {
		fatalErrorf("Error listing import paths : %s", err.Error())
	}
	Log.Info("Found %d imports", len(imports))
	if len(imports) == 0 {
		// No imports, nothing to do
		return
	}

	vendors, err := project.GetVendors()
	if err != nil {
		fatalErrorf("Error listing vendors : %s", err.Error())
	}
	Log.Info("Found %d vendors", len(vendors))
	Log.Debug("Vendors list : %# v", vendors)

	errorCount := 0
	for vendorImportPath := range vendors {
		if !imports.Has(vendorImportPath) {
			err = project.RemoveVendor(vendorImportPath)
			if err != nil {
				errorCount++
				Log.Error("Error removing vendor %q : %s", vendorImportPath, err.Error())
			} else {
				fmt.Println(vendorImportPath)
			}
		}
	}

	if errorCount > 0 {
		os.Exit(errorCount)
	}
}
