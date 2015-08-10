package gpm

import (
	"io/ioutil"
	"testing"

	"os"
	"path"

	"github.com/stretchr/testify/assert"
	"github.com/tsuru/commandmocker"
)

func newProjectNoVCLProject(t *testing.T) (project *ProjectNoVCL) {
	tmpDirPath, err := ioutil.TempDir("", "gpm-project_no_vcl_test")
	assert.Nil(t, err, "Failed creating a temp dir for ProjectNoVCL tests.")

	project = NewProjectNoVCL(tmpDirPath)
	assert.Equal(t, tmpDirPath, project.GetBaseDir(), "Fresh ProjectNoVCL from temp dir has a wrong base dir.")

	return
}

func TestProjectNoVCL_GetVendors(t *testing.T) {
	testCases := map[string]struct {
		prepare             func(t *testing.T, project *ProjectNoVCL, failMessageAndArgs ...interface{}) (tearDown func())
		expectedImportPaths map[string]LocalProject
	}{
		"none": {
			prepare:             func(_ *testing.T, _ *ProjectNoVCL, _ ...interface{}) func() { return func() {} },
			expectedImportPaths: nil,
		},
		"1 NoVCL vendor": {
			prepare: func(t *testing.T, project *ProjectNoVCL, failMessageAndArgs ...interface{}) (tearDown func()) {
				targetDirPath := path.Join(project.GetBaseDir(), "vendor", "test1") + "/"
				err := os.MkdirAll(targetDirPath, os.ModeDir|os.ModePerm)
				assert.Nil(t, err, failMessageAndArgs...)

				targetFilePath := path.Join(targetDirPath, "whatever.go")
				file, err := os.Create(targetFilePath)
				assert.Nil(t, err, failMessageAndArgs...)

				_, err = file.WriteString("package test1\n")
				assert.Nil(t, err, failMessageAndArgs...)

				tmpBinPath, err := commandmocker.Error("git", "fatal: Not a git repository (or any of the parent directories): /fake/path/whatever/test1/", 128)
				assert.Nil(t, err, failMessageAndArgs...)

				tearDown = func() {
					commandmocker.Remove(tmpBinPath)
				}
				return
			},
			expectedImportPaths: map[string]LocalProject{
				"test1": &ProjectNoVCL{},
			},
		},
		"1 Git vendor": {
			prepare: func(t *testing.T, project *ProjectNoVCL, failMessageAndArgs ...interface{}) (tearDown func()) {
				targetDirPath := path.Join(project.GetBaseDir(), "vendor", "test1") + "/"
				err := os.MkdirAll(targetDirPath, os.ModeDir|os.ModePerm)
				assert.Nil(t, err, failMessageAndArgs...)

				targetFilePath := path.Join(targetDirPath, "whatever.go")
				file, err := os.Create(targetFilePath)
				assert.Nil(t, err, failMessageAndArgs...)

				_, err = file.WriteString("package test1\n")
				assert.Nil(t, err, failMessageAndArgs...)

				tmpBinPath, err := commandmocker.Add("git", targetDirPath+"\n")
				assert.Nil(t, err, failMessageAndArgs...)

				tearDown = func() {
					commandmocker.Remove(tmpBinPath)
				}
				return
			},
			expectedImportPaths: map[string]LocalProject{
				"test1": &localGitProject{},
			},
		},
	}

	tearDowns := make(map[string]func(), len(testCases))
	defer func() {
		for _, tearDown := range tearDowns {
			tearDown()
		}
	}()

testCasesLoop:
	for caseName, testCase := range testCases {
		failMessageAndArgs := []interface{}{"Test case %q failed.", caseName}
		project := newProjectNoVCLProject(t)
		defer os.RemoveAll(project.GetBaseDir())

		tearDowns[caseName] = testCase.prepare(t, project, failMessageAndArgs...)

		vendors, err := project.GetVendors()
		if !assert.Nil(t, err, failMessageAndArgs...) {
			continue testCasesLoop
		}

		if !assert.Len(t, vendors, len(testCase.expectedImportPaths), failMessageAndArgs...) {
			continue testCasesLoop
		}

	vendorsLoop:
		for importPath, projectType := range testCase.expectedImportPaths {
			vendor, exists := vendors[importPath]

			specificFailMessageAndArgs := make([]interface{}, len(failMessageAndArgs)+1)
			copy(specificFailMessageAndArgs, failMessageAndArgs)
			specificFailMessageAndArgs[0] = specificFailMessageAndArgs[0].(string) + " (vendor %q)"
			specificFailMessageAndArgs[len(failMessageAndArgs)] = importPath

			if !assert.True(t, exists, specificFailMessageAndArgs...) {
				continue vendorsLoop
			}

			if !assert.NotNil(t, vendor, specificFailMessageAndArgs...) {
				continue vendorsLoop
			}

			if !assert.IsType(t, projectType, vendor.GetProject(), specificFailMessageAndArgs...) {
				continue vendorsLoop
			}
		}

		tearDowns[caseName]()
		delete(tearDowns, caseName)
	}

}
