package main

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/hectorj/godm"
	logging "github.com/op/go-logging"
)

var Log = logging.MustGetLogger("godm")

var submodulesRegexp = regexp.MustCompile(`\[submodule "vendor/([^"]+)"\]`)

func main() {
	checkGo15VendorActivated()

	cli.VersionFlag.Name = "version" // We use "v" for verbose
	app := cli.NewApp()
	app.Name = "godm"
	app.Usage = "Package Manager for Go 1.5+"
	app.Authors = []cli.Author{
		{
			Name:  "HectorJ",
			Email: "hector.jusforgues@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "v,verbose",
			Usage: "Verbose output",
		},
	}
	app.EnableBashCompletion = true
	app.Before = func(c *cli.Context) error {
		backend := logging.AddModuleLevel(logging.NewLogBackend(os.Stderr, "", 0))
		if true || c.GlobalBool("verbose") {
			backend.SetLevel(logging.DEBUG, "")
		} else {
			backend.SetLevel(logging.WARNING, "")
		}

		Log.SetBackend(backend)
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "vendor",
			Usage: "Scans imports from Go files and vendor them in the current Git repository. Takes files/directories path(s) as arguments",
			Flags: []cli.Flag{
			// Removed for now. Here is the current behavior :
			// - For Git projects : we start from the root dir and do scan recursively all sub-packages (except vendors)
			// - For "no-vcl" projects : we only scan the current dir we're in
			//				cli.BoolFlag{
			//					Name:  "r,recursive",
			//					Usage: "Scan dirs recursively for sub-packages",
			//				},

			// Removed for now. Here is the current behavior :
			// We exclude "vendor" directories
			//				cli.StringSliceFlag{
			//					Name:  "e,exclude",
			//					Usage: "Files/Directories names to exclude from scanning, as regexp.",
			//				},
			},
			Action: vendor,
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "Unvendors an import path. Takes a single import path as argument",
			Flags:   []cli.Flag{
			// Removed for now, as there is no confirmation asked anywhere
			//				cli.BoolFlag{
			//					Name:  "y,yes",
			//					Usage: "Remove the submodule without asking any confirmation",
			//				},
			},
			Action: remove,
			BashComplete: func(c *cli.Context) {
				if len(c.Args()) > 0 {
					return
				}
				project, err := godm.NewLocalProject(path.Dir(os.Args[0]))
				if err != nil {
					return
				}

				vendors, err := project.GetVendors()

				if err != nil || vendors == nil {
					return
				}

				for importPath := range vendors {
					fmt.Println(importPath)
				}
			},
		},
	}
	app.RunAndExitOnError()
}

func errorf(format string, args ...interface{}) {
	Log.Error(format, args...)
}

func fatalErrorf(format string, args ...interface{}) {
	errorf(format, args...)
	os.Exit(1)
}

func checkGo15VendorActivated() {
	if os.Getenv("GO15VENDOREXPERIMENT") != "1" {
		Log.Warning("Warning : GO15VENDOREXPERIMENT is not activated.\ngodm relies entirely on that vendoring feature\nTo activate it, run `export GO15VENDOREXPERIMENT=1`")
	}
}
