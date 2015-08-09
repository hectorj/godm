package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func gitAddSubmodule(repoDir, remoteURI, targetPath string) ([]byte, error) {
	cmd := exec.Command("git", "submodule", "add", "-f", remoteURI, targetPath)
	cmd.Dir = repoDir
	return cmd.CombinedOutput()
}

func gitCheckoutCommit(repoDir, commitHash string) ([]byte, error) {
	cmd := exec.Command("git", "checkout", commitHash)
	cmd.Dir = repoDir
	return cmd.CombinedOutput()
}

var remoteExtractRegexp = regexp.MustCompile(`^([^\s]+)\s+([^\s]+) \(fetch\)`)

func gitGetRemoteURI(repoDir string, allowLocal bool) (string, error) {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	matches := remoteExtractRegexp.FindStringSubmatch(string(output))
	if matches == nil {
		if allowLocal {
			// @TODO : Maybe. Vendoring local repo doesn't actually sound like a good idea, gotta see if there is
			// some real usecases
			panic("Getting local repo URI : not implemented yet")
		}
		err = fmt.Errorf("Could not extract remote URL from %q", repoDir)
		return "", err
	}
	return matches[2], nil
}

func gitGetCurrentCommitHash(repoDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	cmd.Dir = repoDir

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(output), "\n"), nil
}
