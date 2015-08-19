package git

var Service GitService = NewGitService()

type GitService interface {
	Clone(targetPath, remoteURI string) error
	AddSubmodule(repoDir, remoteURI, targetPath string) error
	RemoveSubmodule(repoDir, targetPath string) error
	CheckoutCommit(repoDir, commitHash string) error
	GetRemoteURI(repoDir string) (string, error)
	GetCurrentCommitHash(repoDir string) (string, error)
	GetRootDir(dir string) (string, error)
	InitRepo(repoDir string) error
	InitSubmodules(repoDir string) error
	UpdateSubmodules(repoDir string) error
}

func NewGitService() GitService {
	return gitService{}
}

type gitService struct{}
