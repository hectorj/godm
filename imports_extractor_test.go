package gpm

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractImports(t *testing.T) {
	testCases := map[string]struct {
		src    string
		result Set
	}{
		"empty": {
			src:    "package main",
			result: NewSet(),
		},
		"simple": {
			src: `
			package main

			import "test"
			`,
			result: NewSet("test"),
		},
		"multiple": {
			src: `
			package main

			import (
				"test"
				"test2"
			)
			`,
			result: NewSet("test", "test2"),
		},
		"duplicates": {
			src: `
			package main

			import (
				"test"
				"test"
			)
			`,
			result: NewSet("test"),
		},
		"alias": {
			src: `
			package main

			import (
				. "test"
				test "test2"
			)
			`,
			result: NewSet("test", "test2"),
		},
	}

	for caseName, testCase := range testCases {
		errorMessage := fmt.Sprintf("Test case %q failed", caseName)

		file, err := ioutil.TempFile("", "gpm-extractImports-test-"+caseName)
		assert.Nil(t, err, errorMessage)
		defer os.Remove(file.Name())

		io.WriteString(file, testCase.src)

		imports, err := extractImports(NewSet(file.Name()))
		assert.Equal(t, testCase.result, imports, errorMessage)
	}
}
