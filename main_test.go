package main

import (
	"fmt"
	"os"
	"testing"
)

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
	//os.RemoveAll(testRepoDir)
	//

	os.Exit(statusCode)
}
