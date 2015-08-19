package main

import (
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/hectorj/godm"
)

func remove(c *cli.Context) {
	if len(c.Args()) != 1 {
		cli.ShowCommandHelp(c, "remove")
		fatalErrorf("Expected exactly 1 argument instead of %d", len(c.Args()))
	}

	project, err := godm.NewLocalProject(path.Dir(os.Args[0]), "")
	if err != nil {
		fatalErrorf("Error building the current project : %s", err.Error())
	}
	Log.Debug("Project's Type : %T", project, project.GetBaseDir())
	Log.Debug("Project's Base Dir : %s", project.GetBaseDir())

	importPath := c.Args().First()

	err = project.RemoveVendor(importPath)

	if err == godm.ErrUnknownVendor {
		Log.Warning("Import path %q cannot be removed from vendors because it is not vendored in current project", importPath)
		return
	} else if err != nil {
		fatalErrorf("Error removing import path %q from vendors : %s", importPath, err.Error())
	}
	Log.Notice("%q removed from vendors", importPath)
}
