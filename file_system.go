package gpm

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	shutil "github.com/termie/go-shutil"
)

func listGoFiles(currentPath string, recursive bool, excludes []*regexp.Regexp, firstLevel bool) (Set, error) {
	for _, excludeRegexp := range excludes {
		if excludeRegexp.MatchString(currentPath) {
			return nil, nil
		}
	}

	file, err := os.Stat(currentPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if !file.IsDir() {
		if strings.HasSuffix(file.Name(), ".go") {
			return NewSet(currentPath), nil
		} else {
			return nil, nil
		}
	}

	if recursive || firstLevel {
		children, err := ioutil.ReadDir(currentPath)
		if err != nil {
			return nil, err
		}
		set := NewSet()
		for _, child := range children {
			childSet, err := listGoFiles(path.Join(currentPath, child.Name()), recursive, excludes, false)
			if err != nil {
				return nil, err
			}
			set.AddSet(childSet)
		}
		return set, nil
	}
	return nil, nil
}

var copyTreeOptions = &shutil.CopyTreeOptions{
	Symlinks:               true,
	IgnoreDanglingSymlinks: true,
}

func CopyDir(src, target string) error {
	return shutil.CopyTree(src, target, copyTreeOptions)
}
