package main

import (
	"os"

	"flag"
	"fmt"
)

func main() {
	checkGo15VendorActivated()

	app := app{}
	if err := app.parseCommandLineArguments(); err != nil {
		errorf("Error parsing CLI arguments : %s", err.Error())
		flag.Usage() // @TODO : real usage documentation
		os.Exit(1)
	}

	switch app.action {
	case actionVendor:
		imports, err := extractImports(app.vendorParameters.files)
		if err != nil {
			fatalErrorf("Error extracting imports from files : %s", err.Error())
		}
		for _, importPath := range imports {
			fmt.Print(importPath, " : ")
			err, ok := vendorImport(app.currentDir, importPath)
			if err != nil {
				fmt.Print("Failed (", err.Error(), ")")
			} else if !ok {
				fmt.Print("Skipped")
			} else {
				fmt.Print("OK")
			}
			fmt.Print("\n")
		}
	case actionRemove:
		err := removeImport(app.currentDir, app.removeParameters.importPath, app.removeParameters.preApproved)
		if err != nil {
			fatalErrorf("Error removing import : %s", err.Error())
		}
	default:
		panic("Unknown action")
	}
	os.Exit(0)
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
