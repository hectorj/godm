package main

import (
	"os"
	"path"

	"os/exec"

	"fmt"
	"regexp"
	"strings"

	"runtime"

	"golang.org/x/tools/go/vcs"
)

var remoteExtractRegexp = regexp.MustCompile(`^([^\s]+)\s+([^\s]+) \(fetch\)`)

func vendor(dir string, importPath string) (err error, vendored bool) {
	gopath := os.Getenv("GOPATH")
	goroot := runtime.GOROOT()
	pkgDir := path.Join(gopath, "src", importPath)
	if _, err = os.Stat(pkgDir); err != nil {
		if os.IsNotExist(err) {
			pkgDir = path.Join(goroot, "src", importPath)
			if _, err = os.Stat(pkgDir); err != nil {
				if os.IsNotExist(err) {
					err = fmt.Errorf("Import path %q not found in GOPATH %q nor in GOROOT %q", importPath, gopath, goroot)
				} else {
					// It's a standard library package, skipping
					err = nil
				}
			}
		}
		return
	}

	if strings.HasPrefix(pkgDir, dir) {
		// It is a sub-package of the current one, skipping
		return
	}

	var (
		importVCS  *vcs.Cmd
		importRoot string
	)
	importVCS, importRoot, err = vcs.FromDir(pkgDir, gopath)
	if err != nil {
		return
	}
	absoluteImportRoot := path.Join(gopath, importRoot)

	var commitHash, remoteUrl string
	var output []byte
	switch importVCS.Name {
	case "Git":
		cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
		cmd.Dir = absoluteImportRoot

		output, err = cmd.Output()
		if err != nil {
			return
		}
		commitHash = strings.Trim(string(output), "\n")

		cmd = exec.Command("git", "remote", "-v")
		cmd.Dir = absoluteImportRoot
		output, err = cmd.Output()
		if err != nil {
			return
		}
		matches := remoteExtractRegexp.FindStringSubmatch(string(output))
		// @TODO : add an option to allow local repositories ?
		if matches == nil {
			err = fmt.Errorf("Could not extract remote URL from %q", absoluteImportRoot)
			return
		}
		remoteUrl = matches[2]

	default:
		// @TODO : mercurial subrepositories (if `go get` supports them)
		err = fmt.Errorf("Unsupported VCS for dir %q : %q", pkgDir, importVCS.Name)
		return
	}

	var (
		mainVCS  *vcs.Cmd
		mainRoot string
	)
	mainVCS, mainRoot, err = vcs.FromDir(dir, gopath)
	if err != nil {
		return
	}
	absoluteMainRoot := path.Join(gopath, mainRoot)

	targetPath := path.Join(absoluteMainRoot, "vendor", importRoot[4:])
	if _, err = os.Stat(targetPath); err == nil {
		//err = fmt.Errorf("%q already exists", targetPath)
		// Already exists, skipping
		return
	}
	err = nil

	switch mainVCS.Name {
	case "Git":
		cmd := exec.Command("git", "submodule", "add", "-f", remoteUrl, targetPath)
		cmd.Dir = absoluteMainRoot
		output, err = cmd.CombinedOutput()
		if err != nil {
			// @TODO : proper logging with verbose option
			//fmt.Println(string(output))
			return
		}

		cmd = exec.Command("git", "checkout", commitHash)
		cmd.Dir = targetPath
		output, err = cmd.CombinedOutput()
		if err != nil {
			//fmt.Println(string(output))
			return
		}

	default:
		err = fmt.Errorf("Unsupported VCS for dir %q : %q", dir, mainVCS.Name)
		return
	}
	vendored = true
	return
}
