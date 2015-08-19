package git

import (
	"path"
	"strings"

	"fmt"

	"os"

	"github.com/dropbox/godropbox/errors"
)

type GitRepoStub struct {
	RemoteURI         string
	CurrentCommitHash string
	Submodules        map[string]*GitRepoStub
}

func NewGitRepoStub() *GitRepoStub {
	return &GitRepoStub{
		Submodules: make(map[string]*GitRepoStub),
	}
}

type GitStub struct {
	Repos map[string]*GitRepoStub
}

func NewGitStub() *GitStub {
	return &GitStub{
		Repos: make(map[string]*GitRepoStub),
	}
}

func (self *GitStub) Clone(targetPath, remoteURI string) error {
	panic("Not implemented yet!")
	return nil
}
func (self *GitStub) AddSubmodule(repoDir, remoteURI, targetPath string) error {
	for repoPath, repo := range self.Repos {
		if path.Clean(repoPath) == path.Clean(repoDir) {
			absolutePath := path.Clean(path.Join(repoDir, targetPath))
			if _, exists := self.Repos[absolutePath]; exists {
				return errors.New("Repo already exists")
			}
			_, err := os.Stat(absolutePath)
			if !os.IsNotExist(err) {
				return errors.New("Dir already exists")
			}
			os.MkdirAll(absolutePath, os.ModeDir|os.ModePerm)
			repo.Submodules[targetPath] = NewGitRepoStub()
			repo.Submodules[targetPath].RemoteURI = remoteURI
			self.Repos[absolutePath] = repo.Submodules[targetPath]
			return nil
		}
	}
	return errors.New("repo not found")
}
func (self *GitStub) RemoveSubmodule(repoDir, targetPath string) error {
	for repoPath, repo := range self.Repos {
		if path.Clean(repoPath) == path.Clean(repoDir) {
			for submodulePath, _ := range repo.Submodules {
				if path.Clean(submodulePath) == path.Clean(targetPath) {
					absolutePath := path.Clean(path.Join(repoPath, submodulePath))
					os.RemoveAll(absolutePath)
					delete(repo.Submodules, submodulePath)
					delete(self.Repos, absolutePath)
					return nil
				}
			}
			return errors.New("submodule not found")
		}
	}
	return errors.New("repo not found")
}
func (self *GitStub) CheckoutCommit(repoDir, commitHash string) error {
	for repoPath, repo := range self.Repos {
		fmt.Println(repoDir, repoPath)
		if path.Clean(repoPath) == path.Clean(repoDir) {
			repo.CurrentCommitHash = commitHash
			return nil
		}
	}
	fmt.Println("\n")
	return errors.New("repo not found")
}
func (self *GitStub) GetRemoteURI(repoDir string) (string, error) {
	for repoPath, repo := range self.Repos {
		if path.Clean(repoPath) == path.Clean(repoDir) {
			if repo.RemoteURI == "" {
				return "", ErrNoRemote
			}
			return repo.RemoteURI, nil
		}
	}
	return "", errors.New("repo not found")
}
func (self *GitStub) GetCurrentCommitHash(repoDir string) (string, error) {
	for repoPath, repo := range self.Repos {
		if path.Clean(repoPath) == path.Clean(repoDir) {
			return repo.CurrentCommitHash, nil
		}
	}
	return "", errors.New("repo not found")
}
func (self *GitStub) GetRootDir(dir string) (string, error) {
	var longestPath string
	for repoPath := range self.Repos {
		cleanPath := path.Clean(repoPath)
		if strings.HasPrefix(dir, cleanPath) {
			if len(cleanPath) > len(longestPath) {
				longestPath = cleanPath
			}

		}
	}
	if longestPath != "" {
		return longestPath, nil
	}
	return "", ErrNotAGitRepository
}
func (self *GitStub) InitRepo(repoDir string) error {
	panic("Not implemented yet!")
	return nil
}
func (self *GitStub) InitSubmodules(repoDir string) error {
	panic("Not implemented yet!")
	return nil
}
func (self *GitStub) UpdateSubmodules(repoDir string) error {
	panic("Not implemented yet!")
	return nil
}
