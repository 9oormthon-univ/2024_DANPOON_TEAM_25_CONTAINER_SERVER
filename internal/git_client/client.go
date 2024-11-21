package gitclient

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

const REPO_URL = "https://github.com/9oormthon-univ/2024_DANPOON_TEAM_25_BE.git"
const CLONE_DIR = "./clone_repo"
const TARGET_DIR = "IDE"
const DEPLOYMENT_FILENAME = "deployment.yaml"
const SERVICE_FILENAME = "service.yaml"
const COMMIT_MESSAGE = "bot: add IDE by userID:%s on courseID: %s"
const BRANCH_NAME = "main"

type Gitclient struct {
	Auth *githttp.BasicAuth
}

func NewAuth() *githttp.BasicAuth {
	auth := &githttp.BasicAuth{
		Username: os.Getenv("GITHUB_USERNAME"),
		Password: os.Getenv("GITHUB_CREDENTIALS"),
	}
	return auth
}

func NewGitClient() *Gitclient {
	return &Gitclient{Auth: NewAuth()}
}

func (g *Gitclient) ModifyRepository() error {
	repo, err := git.PlainClone(CLONE_DIR, false, &git.CloneOptions{
		URL:      REPO_URL,
		Progress: os.Stdout,
		Auth:     g.Auth,
	})
	if err != nil {
		return err
	}
	_, err = repo.Worktree()
	if err != nil {
		return err
	}

	targetPath := filepath.Join(CLONE_DIR, TARGET_DIR)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		log.Fatalf("Fail to create directory: %v", err)
	}

	//TODO: deployment.yaml 생성
	//TODO: service.yaml 생성
	//TODO: commit, push
	return nil
}
