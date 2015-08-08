package main

import (
	"go/token"
	"os"

	"fmt"
	"go/ast"

	"go/parser"
	"path"
	"path/filepath"
)

func main() {
	fs := token.NewFileSet()
	v := &visitor{
		ImportPathsMap: make(map[string]struct{}),
	}
	for _, fileName := range os.Args[1:] {
		fileInfo, err := os.Stat(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error with file %q : %q\n", fileName, err.Error())
			os.Exit(1)
		}
		if fileInfo.IsDir() {
			// @TODO : allow directory scanning (with an option for recursive scanning)
			fmt.Fprintf(os.Stderr, "%q is a dir (not supported yet)", fileName)
			os.Exit(1)
		}
		file, err := parser.ParseFile(fs, fileName, nil, parser.ImportsOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error with file %q : %q\n", fileName, err.Error())
			os.Exit(1)
		}
		ast.Walk(v, file)
	}

	currentDir, err := filepath.Abs(path.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting absolute current path : %q", err.Error())
		os.Exit(1)
	}

	for importPath, _ := range v.ImportPathsMap {
		fmt.Print(importPath, " : ")
		err, ok := vendor(currentDir, importPath)
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
