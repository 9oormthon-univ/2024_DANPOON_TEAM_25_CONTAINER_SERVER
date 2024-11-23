package gitclient

import (
	"bytes"
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
	"golang.org/x/crypto/ssh"
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

func (g *Gitclient) ModifyRepository(courseID, studentID string) error {
	imageTag := fmt.Sprintf("course%s", courseID)
	encodedImageTag := base64.RawURLEncoding.EncodeToString([]byte(imageTag))
	ideTag := fmt.Sprintf("user%scourse%s", studentID, courseID)
	repo, err := git.PlainClone(util.GetPath(CLONE_DIR), false, &git.CloneOptions{
		URL:      REPO_URL,
		Progress: os.Stdout,
		Auth:     g.Auth,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	appFilePath := filepath.Join(util.GetPath(CLONE_DIR), fmt.Sprintf("ide-%s", ideTag))
	err = os.Mkdir(appFilePath, os.FileMode(0777))
	if err != nil {
		return err
	}

	deploymentYaml := fmt.Sprintf(DEPLOYMENT_MANIFEST, ideTag, ideTag, ideTag, ideTag, ideTag, encodedImageTag)
	err = os.WriteFile(fmt.Sprintf("%s/deployment.yaml", appFilePath), []byte(deploymentYaml), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Deployment 파일 생성 완료: %s\n", appFilePath)

	serviceYaml := fmt.Sprintf(SERVICE_MANIFEST, ideTag, ideTag)
	err = os.WriteFile(fmt.Sprintf("%s/service.yaml", appFilePath), []byte(serviceYaml), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Service 파일 생성 완료: %s\n", appFilePath)

	//TODO: Ingress Route 추가
	ingressRouteYaml := fmt.Sprintf(INGRESS_ROUTE_MANIFEST, ideTag, fmt.Sprintf("`%s.flakeide.com`", encodedImageTag), ideTag)
	ingressRoutePath := filepath.Join(util.GetPath(CLONE_DIR), "traefik")
	file, err := os.OpenFile(fmt.Sprintf("%s/ingress-routes.yaml", ingressRoutePath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(ingressRouteYaml)
	if err != nil {
		return err
	}
	applicationYaml := fmt.Sprintf(APPLICATION_MANIFEST, studentID, courseID, studentID, courseID)
	applicationFile, err := os.OpenFile(fmt.Sprintf("%s/application.yaml", util.GetPath(CLONE_DIR)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = applicationFile.WriteString(applicationYaml)
	if err != nil {
		return err
	}

	//TODO: commit,push
	_, err = worktree.Add(".")
	if err != nil {
		log.Fatalf("파일 추가 실패: %v", err)
	}
	fmt.Println("파일 Git에 추가 완료.")
	commit, err := worktree.Commit(fmt.Sprintf("bot: create ide-%s space", ideTag), &git.CommitOptions{
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

	server := "213.190.4.144:22"
	username := "root"
	password := "Flakeide123!" // 키를 사용하는 경우는 비워두세요.

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	client, err := ssh.Dial("tcp", server, config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	command := "cd 2024_DANPOON_TEAM_25_MANIFEST && git pull origin main && kubectl apply -f application.yaml"
	if err := session.Run(command); err != nil {
		log.Fatalf("Failed to run command: %s", err)
	}

	fmt.Printf("Output:\n%s", stdout.String())
	if stderr.Len() > 0 {
		fmt.Printf("Error Output:\n%s", stderr.String())
	}
	return nil
}
