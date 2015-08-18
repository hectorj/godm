package godm

import (
	"io/ioutil"
	"testing"

	"os"

	"path"

	"github.com/stretchr/testify/assert"
)

func TestNewProjectFromImportPath(t *testing.T) {
	tmpDirPath, err := ioutil.TempDir("", "godm-project_test")
	assert.Nil(t, err)

	defer os.Setenv("GOPATH", os.Getenv("GOPATH"))
	os.Setenv("GOPATH", tmpDirPath)

	testImportPath := "test1"
	testSubimportPath := testImportPath + "/test2"

	testPath := path.Join(tmpDirPath, "src", testImportPath)

	subtestPath := path.Join(tmpDirPath, "src", testSubimportPath)

	err = os.MkdirAll(subtestPath, os.ModeDir|os.ModePerm)
	assert.Nil(t, err)

	//// Case 1
	subproject, canonicalImportPath, err := NewProjectFromImportPath(testSubimportPath)
	assert.Nil(t, err)

	assert.IsType(t, new(ProjectNoVCL), subproject)

	localSubproject := subproject.(LocalProject)
	assert.Equal(t, testSubimportPath, canonicalImportPath)
	assert.Equal(t, subtestPath, localSubproject.GetBaseDir())

	//// Case 2
	project, canonicalImportPath, err := NewProjectFromImportPath(testImportPath)
	assert.Nil(t, err)

	assert.IsType(t, new(ProjectNoVCL), project)

	localProject := project.(LocalProject)
	assert.Equal(t, testImportPath, canonicalImportPath)
	assert.Equal(t, testPath, localProject.GetBaseDir())

	//// Case 3
	_, canonicalImportPath, err = NewProjectFromImportPath("fmt")

	assert.NotNil(t, err)
	assert.Equal(t, ErrStandardLibrary, err)

	assert.Equal(t, "fmt", canonicalImportPath)

	//// Case 4
	os.MkdirAll(path.Join(tmpDirPath, "src", "ThatPathIsNotReadable"), os.ModeDir)
	_, _, err = NewProjectFromImportPath("ThatPathIsNotReadable")

	assert.NotNil(t, err)
	assert.False(t, os.IsNotExist(err))
}
