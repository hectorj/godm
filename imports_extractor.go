package godm

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type importsExtractor struct {
	ImportPaths Set
}

func extractImports(fileNames Set) (Set, error) {
	if len(fileNames) == 0 {
		return nil, nil
	}

	fs := token.NewFileSet()
	v := &importsExtractor{
		ImportPaths: NewSet(),
	}
	for fileName := range fileNames {
		file, err := parser.ParseFile(fs, fileName, nil, parser.ImportsOnly)
		if err != nil {
			return nil, fmt.Errorf("Error with file %q : %q\n", fileName, err.Error())
		}
		ast.Walk(v, file)
	}
	return v.ImportPaths, nil
}

func (self *importsExtractor) Visit(node ast.Node) ast.Visitor {
	if importDeclaration, ok := node.(*ast.ImportSpec); ok {
		self.ImportPaths.Add(strings.Trim(importDeclaration.Path.Value, `"`))
	}
	return self
}
