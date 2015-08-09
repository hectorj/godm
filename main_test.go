package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

var testRepoDir string

const (
	testRepoZip                   = "./gpm-test-repo.zip"
	testRepoRemote                = "https://github.com/hectorj/gpm-test-repo.git"
	testRepoCommitHash            = "9bc1c419ff8b07737b880e2bacd2f7d029c91b69"
	testRepoExistingSubmodulePath = "submodule-to-remove"
	testRepoExistingSubdirPath    = "subdir"
	testRepoCommitMessage         = "Cats."
)

func TestMain(m *testing.M) {
	// Setup
	tmpDir, err := ioutil.TempDir("", "gpm-tests")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	err = unzip(testRepoZip, tmpDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(tmpDir)
	testRepoDir = path.Join(tmpDir, "gpm-test-repo")
	//

	statusCode := m.Run()

	// Teardown
	os.RemoveAll(tmpDir)
	//

	os.Exit(statusCode)
}

// from http://stackoverflow.com/a/24792688/1685538
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
