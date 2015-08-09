package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type visitor struct {
	ImportPathsMap map[string]struct{}
	ImportPaths    []string
}

func extractImports(fileNames []string) ([]string, error) {
	fs := token.NewFileSet()
	v := &visitor{
		ImportPathsMap: make(map[string]struct{}),
	}
	for _, fileName := range fileNames {
		file, err := parser.ParseFile(fs, fileName, nil, parser.ImportsOnly)
		if err != nil {
			return nil, fmt.Errorf("Error with file %q : %q\n", fileName, err.Error())
		}
		ast.Walk(v, file)
	}
	return v.ImportPaths, nil
}

func (self *visitor) Visit(node ast.Node) ast.Visitor {
	if importDeclaration, ok := node.(*ast.ImportSpec); ok {
		importPath := strings.Trim(importDeclaration.Path.Value, `"`)
		if _, exists := self.ImportPathsMap[importPath]; !exists {
			self.ImportPathsMap[importPath] = struct{}{}
			self.ImportPaths = append(self.ImportPaths, importPath)
		}
	}
	return self
}
