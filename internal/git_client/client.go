package gitclient

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER/internal/util"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

const REPO_URL = "https://github.com/9oormthon-univ/2024_DANPOON_TEAM_25_MANIFEST.git"
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

func (g *Gitclient) ModifyRepository(key string) error {
	repo, err := git.PlainClone(util.GetPath(CLONE_DIR), false, &git.CloneOptions{
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
	appFilePath := filepath.Join(util.GetPath(CLONE_DIR), fmt.Sprintf("ide-%s", key))
	err = os.Mkdir(appFilePath, os.FileMode(0777))
	if err != nil {
		return err
	}

	deploymentYaml := fmt.Sprintf(DEPLOYMENT_MANIFEST, key, key, key, key, key, key)
	err = os.WriteFile(fmt.Sprintf("%s/deployment.yaml", appFilePath), []byte(deploymentYaml), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Deployment 파일 생성 완료: %s\n", appFilePath)
	servieYaml := fmt.Sprintf(SERVICE_MANIFEST, key, key)
	err = os.WriteFile(fmt.Sprintf("%s/service.yaml", appFilePath), []byte(servieYaml), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Service 파일 생성 완료: %s\n", appFilePath)

	err = os.RemoveAll(util.GetPath(CLONE_DIR))
	if err != nil {
		return err
	}

	//TODO: commit, push
	return nil
}
