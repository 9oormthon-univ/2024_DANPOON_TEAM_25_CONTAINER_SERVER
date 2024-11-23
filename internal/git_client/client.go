package gitclient

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER/internal/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}
	decodedKey := string(decodedBytes)
	log.Println(decodedKey)
	appFilePath := filepath.Join(util.GetPath(CLONE_DIR), fmt.Sprintf("ide.%s", decodedKey))
	err = os.Mkdir(appFilePath, os.FileMode(0777))
	if err != nil {
		return err
	}

	deploymentYaml := fmt.Sprintf(DEPLOYMENT_MANIFEST, decodedKey, decodedKey, decodedKey, decodedKey, decodedKey, key)
	err = os.WriteFile(fmt.Sprintf("%s/deployment.yaml", appFilePath), []byte(deploymentYaml), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Deployment 파일 생성 완료: %s\n", appFilePath)

	serviceYaml := fmt.Sprintf(SERVICE_MANIFEST, decodedKey, decodedKey)
	err = os.WriteFile(fmt.Sprintf("%s/service.yaml", appFilePath), []byte(serviceYaml), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Service 파일 생성 완료: %s\n", appFilePath)

	//TODO: Ingress Route 추가
	ingressRouteYaml := fmt.Sprintf(INGRESS_ROUTE_MANIFEST, decodedKey, fmt.Sprintf("`%s.flakeide.com`", key), decodedKey)
	ingressRoutePath := filepath.Join(util.GetPath(CLONE_DIR), "traefik")
	file, err := os.OpenFile(fmt.Sprintf("%s/ingress-routes.yaml", ingressRoutePath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(ingressRouteYaml)
	if err != nil {
		return err
	}

	//TODO: commit,push
	_, err = worktree.Add(".")
	if err != nil {
		log.Fatalf("파일 추가 실패: %v", err)
	}
	fmt.Println("파일 Git에 추가 완료.")
	commit, err := worktree.Commit(fmt.Sprintf("bot: create ide-%s space", decodedKey), &git.CommitOptions{
		Author: &object.Signature{
			Name:  os.Getenv("GITHUB_NAME"),
			Email: os.Getenv("GITHUB_EMAIL"),
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("commit complete: %v\n", commit.String())
	err = repo.Push(&git.PushOptions{
		Auth: g.Auth,
	})
	if err != nil {
		log.Printf("Push Error: %v", err)
		return err
	}

	err = os.RemoveAll(util.GetPath(CLONE_DIR))
	if err != nil {
		return err
	}
	return nil
}
