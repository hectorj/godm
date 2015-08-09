package main

import (
	"os"

	"fmt"

	"github.com/codegangsta/cli"
)

func main() {
	checkGo15VendorActivated()

	app := cli.NewApp()
	app.Name = "gpm"
	app.Usage = "Package Manager for Go 1.5+"
	app.Authors = []cli.Author{
		{
			Name:  "HectorJ",
			Email: "hector.jusforgues@gmail.com",
		},
	}
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:  "vendor",
			Usage: "Scans imports from Go files and vendor them in the current Git repository. Takes files/directories path(s) as arguments",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "r",
					Usage: "Scan dirs recursively",
				},
			},
			Action: vendor,
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "Unvendors an import path. Takes a single import path as argument",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "y",
					Usage: "Remove the submodule without asking any confirmation",
				},
			},
			Action: remove,
		},
	}
	app.RunAndExitOnError()
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func fatalErrorf(format string, args ...interface{}) {
	errorf(format, args...)
	os.Exit(1)
}

func checkGo15VendorActivated() {
	if os.Getenv("GO15VENDOREXPERIMENT") != "1" {
		fmt.Fprint(os.Stderr, "Warning : GO15VENDOREXPERIMENT is not activated.\ngpm relies entirely on that vendoring feature\nTo activate it, run `export GO15VENDOREXPERIMENT=1`\n")
	}
}
