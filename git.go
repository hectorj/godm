package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var gitCommand string

func init() {
	var err error
	gitCommand, err = exec.LookPath("git")
	if err != nil {
		panic(err)
	}
}

func gitClone(targetPath, remoteURI string) ([]byte, error) {
	cmd := exec.Command(gitCommand, "clone", remoteURI, targetPath)
	return cmd.CombinedOutput()
}

func gitAddSubmodule(repoDir, remoteURI, targetPath string) ([]byte, error) {
	cmd := exec.Command(gitCommand, "submodule", "add", "-f", remoteURI, targetPath)
	cmd.Dir = repoDir
	return cmd.CombinedOutput()
}

func gitRemoveSubmodule(repoDir, targetPath string) ([]byte, error) {
	buf := bytes.Buffer{}

	cmd := exec.Command(gitCommand, "submodule", "deinit", targetPath)
	cmd.Dir = repoDir
	tmp, err := cmd.CombinedOutput()
	if err != nil {
		return tmp, err
	}
	buf.Write(tmp)

	cmd = exec.Command(gitCommand, "rm", "-rf", targetPath)
	cmd.Dir = repoDir
	tmp, err = cmd.CombinedOutput()
	if err != nil {
		return buf.Bytes(), err
	}
	buf.Write(tmp)

	cmd = exec.Command("rm", "-rf", path.Join(".git/modules/", targetPath))
	cmd.Dir = repoDir
	tmp, err = cmd.CombinedOutput()
	if err != nil {
		return buf.Bytes(), err
	}
	buf.Write(tmp)

	return buf.Bytes(), nil
}

func gitCheckoutCommit(repoDir, commitHash string) ([]byte, error) {
	cmd := exec.Command(gitCommand, "checkout", commitHash)
	cmd.Dir = repoDir
	return cmd.CombinedOutput()
}

var remoteExtractRegexp = regexp.MustCompile(`^([^\s]+)\s+([^\s]+) \(fetch\)`)

func gitGetRemoteURI(repoDir string, allowLocal bool) (string, error) {
	cmd := exec.Command(gitCommand, "remote", "-v")
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
	cmd := exec.Command(gitCommand, "rev-parse", "--verify", "HEAD")
	cmd.Dir = repoDir

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(output), "\n"), nil
}

func gitGetRootDir(dir string) (string, error) {
	cmd := exec.Command(gitCommand, "rev-parse", "--show-toplevel")
	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(output), "\n"), nil
}
