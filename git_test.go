package main

import (
	"testing"

	"os"

	"fmt"

	"path"

	"io/ioutil"

	"path/filepath"

	"github.com/stretchr/testify/assert"
)

func init() {

}

func TestMain(m *testing.M) {
	// Setup
	os.RemoveAll(testRepoDir)
	ouput, err := gitClone(testRepoDir, testRepoRemote)
	if err != nil {
		fmt.Fprintln(os.Stderr, string(ouput))
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	ouput, err = gitInitSubmodules(testRepoDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, string(ouput))
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	ouput, err = gitUpdateSubmodules(testRepoDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, string(ouput))
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	//

	statusCode := m.Run()

	// Teardown
	os.RemoveAll(testRepoDir)
	//

	os.Exit(statusCode)
}

const (
	testRepoDir                   = "./gpm-test-repo"
	testRepoRemote                = "https://github.com/hectorj/gpm-test-repo.git"
	testRepoCommitHash            = "9bc1c419ff8b07737b880e2bacd2f7d029c91b69"
	testRepoExistingSubmodulePath = "submodule-to-remove"
	testRepoExistingSubdirPath    = "subdir"
	testRepoCommitMessage         = "Cats."
)

func TestGitGetRemoteURI(t *testing.T) {
	remote, err := gitGetRemoteURI(testRepoDir, false)

	assert.Nil(t, err)
	assert.Equal(t, testRepoRemote, remote)
}

func TestGitGetCurrentCommitHash(t *testing.T) {
	hash, err := gitGetCurrentCommitHash(testRepoDir)

	assert.Nil(t, err)
	assert.Equal(t, testRepoCommitHash, hash)
}

func TestGitAddSubmodule(t *testing.T) {
	targetPath := "inception"

	output, err := gitAddSubmodule(testRepoDir, testRepoRemote, targetPath)

	assert.Nil(t, err)
	assert.Equal(t, "Cloning into 'inception'...\n", string(output))
}

func TestGitCheckoutCommit(t *testing.T) {
	output, err := gitCheckoutCommit(testRepoDir, testRepoCommitHash)

	assert.Nil(t, err)
	assert.Contains(t, string(output), "Note: checking out '"+testRepoCommitHash+"'.\n\nYou are in 'detached HEAD' state. You can look around, make experimental\nchanges and commit them, and you can discard any commits you make in this\nstate without impacting any branches by performing another checkout.\n\nIf you want to create a new branch to retain commits you create, you may\ndo so (now or later) by using -b with the checkout command again. Example:\n\n  git checkout -b new_branch_name\n\nHEAD is now at "+testRepoCommitHash[:7]+"... "+testRepoCommitMessage+"\n")
}

func TestGitRemoveSubmodule(t *testing.T) {
	output, err := gitRemoveSubmodule(testRepoDir, testRepoExistingSubmodulePath)

	assert.Nil(t, err)
	assert.Equal(t, "Cleared directory 'submodule-to-remove'\nSubmodule 'submodule-to-remove' (https://github.com/hectorj/gpm-test-repo.git) unregistered for path 'submodule-to-remove'\nrm 'submodule-to-remove'\n", string(output))

	_, err = os.Stat(path.Join(testRepoDir, testRepoExistingSubmodulePath))
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(path.Join(testRepoDir, ".git/modules/", testRepoExistingSubmodulePath))
	assert.True(t, os.IsNotExist(err))

	content, err := ioutil.ReadFile(path.Join(testRepoDir, ".gitmodules"))
	if err != nil {
		assert.True(t, os.IsNotExist(err))
	} else {
		assert.NotRegexp(t, `\[submodule "`+testRepoExistingSubmodulePath+`"\]`, string(content))
	}

}

func TestGitGetRootDir(t *testing.T) {
	absoluteTestRepoDir, err := filepath.Abs(testRepoDir)
	assert.Nil(t, err)
	root, err := gitGetRootDir(path.Join(testRepoDir, testRepoExistingSubdirPath))
	assert.Nil(t, err)
	assert.Equal(t, absoluteTestRepoDir, root)
}
