package main

import (
	"os"

	"fmt"

	"io/ioutil"

	"path"
	"path/filepath"

	"regexp"

	"github.com/codegangsta/cli"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("gpm")

var submodulesRegexp = regexp.MustCompile(`\[submodule "vendor/([^"]+)"\]`)

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
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Verbose",
		},
	}
	app.EnableBashCompletion = true
	app.Before = func(c *cli.Context) error {
		backend := logging.AddModuleLevel(logging.NewLogBackend(os.Stderr, "", 0))
		if c.Bool("verbose") {
			backend.SetLevel(logging.DEBUG, "")
		} else {
			backend.SetLevel(logging.WARNING, "")
		}

		log.SetBackend(backend)
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "vendor",
			Usage: "Scans imports from Go files and vendor them in the current Git repository. Takes files/directories path(s) as arguments",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "r,recursive",
					Usage: "Scan dirs recursively",
				},
				cli.StringSliceFlag{
					Name:  "e,exclude",
					Usage: "Files/Directories names to exclude from scanning, as regexp.",
				},
			},
			Action: vendor,
			BashComplete: func(c *cli.Context) {
				// This will complete if no args are passed
				fmt.Println(".")
				fileInfos, err := ioutil.ReadDir(".")
				if err != nil {
					return
				}
				for _, fileInfo := range fileInfos {
					fmt.Println(fileInfo.Name())
				}
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "Unvendors an import path. Takes a single import path as argument",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "y,yes",
					Usage: "Remove the submodule without asking any confirmation",
				},
			},
			Action: remove,
			BashComplete: func(c *cli.Context) {
				currentDir, err := filepath.Abs(path.Dir(os.Args[0]))
				if err != nil {
					return
				}
				gitRoot, err := gitGetRootDir(currentDir)
				if err != nil {
					return
				}
				content, err := ioutil.ReadFile(path.Join(gitRoot, ".gitmodules"))
				if err != nil {
					return
				}

				matches := submodulesRegexp.FindAllStringSubmatch(string(content), -1)
				for _, match := range matches {
					fmt.Println(match[1])
				}
			},
		},
	}
	app.RunAndExitOnError()
}

func errorf(format string, args ...interface{}) {
	log.Error(format, args...)
}

func fatalErrorf(format string, args ...interface{}) {
	errorf(format, args...)
	os.Exit(1)
}

func checkGo15VendorActivated() {
	if os.Getenv("GO15VENDOREXPERIMENT") != "1" {
		log.Warning("Warning : GO15VENDOREXPERIMENT is not activated.\ngpm relies entirely on that vendoring feature\nTo activate it, run `export GO15VENDOREXPERIMENT=1`")
	}
}
