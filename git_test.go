package main

import (
	"testing"

	"os"
	"path"

	"fmt"

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
		panic(err)
	}
	// Run
	statusCode := m.Run()
	// Teardown
	os.RemoveAll(testRepoDir)
	// Exit
	os.Exit(statusCode)
}

const (
	testRepoDir        = "./gpm-test-repo"
	testRepoRemote     = "https://github.com/hectorj/gpm-test-repo.git"
	testRepoCommitHash = "9e6e573e7ffa51a77d68e20027195b22ef26b47a"
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
	submoduleCleanup(targetPath)

	output, err := gitAddSubmodule(testRepoDir, testRepoRemote, targetPath)
	assert.Nil(t, err)
	assert.Equal(t, "Cloning into 'inception'...\n", string(output))

	submoduleCleanup(targetPath)
}

func submoduleCleanup(targetPath string) {
	os.RemoveAll(path.Join(testRepoDir, targetPath))
	os.Remove(path.Join(testRepoDir, ".gitmodules"))
}

func TestGitCheckoutCommit(t *testing.T) {
	output, err := gitCheckoutCommit(testRepoDir, testRepoCommitHash)
	assert.Nil(t, err)
	assert.Equal(t, "Note: checking out '9e6e573e7ffa51a77d68e20027195b22ef26b47a'.\n\nYou are in 'detached HEAD' state. You can look around, make experimental\nchanges and commit them, and you can discard any commits you make in this\nstate without impacting any branches by performing another checkout.\n\nIf you want to create a new branch to retain commits you create, you may\ndo so (now or later) by using -b with the checkout command again. Example:\n\n  git checkout -b new_branch_name\n\nHEAD is now at 9e6e573... Initial commit\n", string(output))
}
