package main

import (
	"fmt"
	"testing"

	"io"
	"io/ioutil"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestExtractImports(t *testing.T) {
	testCases := map[string]struct {
		src    string
		result []string
	}{
		"empty": {
			src:    "package main",
			result: nil,
		},
		"simple": {
			src: `
			package main

			import "test"
			`,
			result: []string{
				"test",
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
			result: []string{
				"test",
				"test2",
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
			result: []string{
				"test",
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
			result: []string{
				"test",
				"test2",
			},
		},
	}

	for caseName, testCase := range testCases {
		errorMessage := fmt.Sprintf("Test case %q failed", caseName)

		file, err := ioutil.TempFile("", "gpm-extractImports-test-"+caseName)
		assert.Nil(t, err, errorMessage)
		defer os.Remove(file.Name())

		io.WriteString(file, testCase.src)

		imports, err := extractImports([]string{file.Name()})
		assert.Equal(t, testCase.result, imports, errorMessage)
	}
}
