package main

import (
	"go/ast"
	"strings"
)

type visitor struct {
	ImportPathsMap map[string]struct{}
}

func (self *visitor) Visit(node ast.Node) ast.Visitor {
	if importDeclaration, ok := node.(*ast.ImportSpec); ok {
		self.ImportPathsMap[strings.Trim(importDeclaration.Path.Value, `"`)] = struct{}{}
	}
	return self
}
