package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisitor(t *testing.T) {
	testCases := map[string]struct {
		src    string
		result map[string]struct{}
	}{
		"empty": {
			src:    "package main",
			result: map[string]struct{}{},
		},
		"simple": {
			src: `
			package main

			import "test"
			`,
			result: map[string]struct{}{
				"test": struct{}{},
			},
		},
		"multiple": {
			src: `
			package main

			import (
				"test"
				"test2"
			)
			`,
			result: map[string]struct{}{
				"test":  struct{}{},
				"test2": struct{}{},
			},
		},
		"duplicates": {
			src: `
			package main

			import (
				"test"
				"test"
			)
			`,
			result: map[string]struct{}{
				"test": struct{}{},
			},
		},
		"alias": {
			src: `
			package main

			import (
				. "test"
				test "test2"
			)
			`,
			result: map[string]struct{}{
				"test":  struct{}{},
				"test2": struct{}{},
			},
		},
	}

	for caseName, testCase := range testCases {
		errorMessage := fmt.Sprintf("Test case %q failed", caseName)
		v := &visitor{
			ImportPathsMap: make(map[string]struct{}),
		}

		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, "src.go", testCase.src, parser.ParseComments)
		assert.Nil(t, err, errorMessage)

		ast.Walk(v, f)
		assert.Equal(t, testCase.result, v.ImportPathsMap, errorMessage)
	}
}
