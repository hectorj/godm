package main

import (
	"os"
	"path"

	"strings"

	"github.com/codegangsta/cli"
	"github.com/hectorj/godm"
)

func vendor(c *cli.Context) {
	if len(c.Args()) > 0 {
		cli.ShowCommandHelp(c, "vendor")
		fatalErrorf("Too much arguments")
	}

	project, err := godm.NewLocalProject(path.Dir(os.Args[0]))
	if err != nil {
		fatalErrorf("Error building the current project : %s", err.Error())
	}
	Log.Debug("Project's Type : %T\nProject's Base Dir : %s", project, project.GetBaseDir())

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

importsLoop:
	for importPath := range imports {
		if _, exists := vendors[importPath]; !exists {
			vendorProject, canonicalImportPath, err := godm.NewProjectFromImportPath(importPath)
			if _, exists := vendors[canonicalImportPath]; exists {
				Log.Info("%q is a sub-package of %q which is already vendored", importPath, canonicalImportPath)
				continue importsLoop
			}

			if localVendorProject, ok := vendorProject.(godm.LocalProject); ok {
				if strings.HasPrefix(localVendorProject.GetBaseDir(), project.GetBaseDir()) {
					Log.Debug("%q ignored because it is a sub-package of current project", canonicalImportPath)
					continue importsLoop
				}
			}
			if err != nil {
				if err == godm.ErrStandardLibrary {
					Log.Debug("%q ignored because part of the standard library", canonicalImportPath)
					err = nil
				} else {
					fatalErrorf("Error with import path %q : %s", canonicalImportPath, err.Error())
				}
			} else {
				_, err = project.AddVendor(canonicalImportPath, vendorProject)
				if err != nil {
					fatalErrorf("Error vendoring import path %q : %s", canonicalImportPath, err.Error())
				} else {
					Log.Notice("%q vendored", canonicalImportPath)
				}
			}

		} else {
			Log.Info("%q is already vendored", importPath)
		}
	}
}
