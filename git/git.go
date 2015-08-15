package git

import (
	"bytes"
	"errors"
	"fmt"
	exec2 "os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/hectorj/godm/exec"
)

var gitCommand = "git"

func init() {
	var err error
	_, err = exec2.LookPath("git")
	if err != nil {
		panic(err)
	}
}

func Clone(targetPath, remoteURI string) error {
	return exec.Cmd("", gitCommand, "clone", remoteURI, targetPath).GetError()
}

func AddSubmodule(repoDir, remoteURI, targetPath string) error {
	return exec.Cmd(repoDir, gitCommand, "submodule", "add", "-f", remoteURI, targetPath).GetError()
}

func RemoveSubmodule(repoDir, targetPath string) error {
	result := exec.Cmd(repoDir, gitCommand, "submodule", "deinit", "-f", targetPath)
	if err := result.GetError(); err != nil {
		return err
	}

	result = exec.Cmd(repoDir, gitCommand, "rm", "-rf", targetPath)
	if err := result.GetError(); err != nil {
		return err
	}

	return exec.Cmd(repoDir, "rm", "-rf", path.Join(".git/modules/", targetPath)).GetError()
}

func CheckoutCommit(repoDir, commitHash string) error {
	return exec.Cmd(repoDir, gitCommand, "checkout", commitHash).GetError()
}

var remoteExtractRegexp = regexp.MustCompile(`^([^\s]+)\s+([^\s]+) \(fetch\)`)

var ErrNoRemote = errors.New("No remote found")

func GetRemoteURI(repoDir string) (string, error) {
	result := exec.Cmd(repoDir, gitCommand, "remote", "-v")

	if err := result.GetError(); err != nil {
		return "", err
	}
	if len(bytes.Trim(result.GetStdout(), "\n")) == 0 {
		return "", ErrNoRemote
	}
	matches := remoteExtractRegexp.FindStringSubmatch(string(result.GetStdout()))
	if matches == nil {
		return "", fmt.Errorf("Could not extract remote URL from %q", repoDir)
	}
	return matches[2], nil
}

func GetCurrentCommitHash(repoDir string) (string, error) {
	result := exec.Cmd(repoDir, gitCommand, "rev-parse", "--verify", "HEAD")

	if err := result.GetError(); err != nil {
		return "", err
	}
	return strings.Trim(string(result.GetStdout()), "\n"), nil
}

var ErrNotAGitRepository = errors.New("Not a git repository")

func GetRootDir(dir string) (string, error) {
	result := exec.Cmd(dir, gitCommand, "rev-parse", "--show-toplevel")

	if err := result.GetError(); err != nil {
		if bytes.Contains(result.GetStderr(), []byte("Not a git repository")) {
			return "", ErrNotAGitRepository
		}
		return "", err
	}
	return strings.Trim(string(result.GetStdout()), "\n"), nil
}

func InitSubmodules(repoDir string) error {
	return exec.Cmd(repoDir, gitCommand, "submodule", "init").GetError()
}

func UpdateSubmodules(repoDir string) error {
	return exec.Cmd(repoDir, gitCommand, "submodule", "update").GetError()
}
