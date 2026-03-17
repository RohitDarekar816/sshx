package git

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type UsersFile struct {
	Users []interface{} `json:"users"`
}

func CloneOrInitRepo(repoURL, repoDir string) error {

	// try cloning
	_, err := gogit.PlainClone(repoDir, false, &gogit.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err == nil {
		fmt.Println("Repository cloned successfully")
		return nil
	}

	fmt.Println("Repo appears empty. Initializing...")

	// ensure directory exists
	err = os.MkdirAll(repoDir, 0755)
	if err != nil {
		return err
	}

	// initialize repo
	repo, err := gogit.PlainInit(repoDir, false)
	if err != nil {

		// repo may already exist locally
		repo, err = gogit.PlainOpen(repoDir)
		if err != nil {
			return err
		}
	}

	// configure remote
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})

	if err != nil {
		fmt.Println("Remote may already exist, continuing...")
	}

	usersFile := filepath.Join(repoDir, "users.json")

	initial := map[string]interface{}{
		"users": []interface{}{},
	}

	data, err := json.MarshalIndent(initial, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(usersFile, data, 0644)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add("users.json")
	if err != nil {
		return err
	}

	_, err = w.Commit("Initialize sshx repo", &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "sshx",
			Email: "sshx@local",
			When:  time.Now(),
		},
		Committer: &object.Signature{
			Name:  "sshx",
			Email: "sshx@local",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	err = repo.Push(&gogit.PushOptions{})
	if err != nil {
		fmt.Println("Push may fail if repo already initialized, continuing...")
	}

	fmt.Println("Repository initialized and pushed")

	return nil
}
