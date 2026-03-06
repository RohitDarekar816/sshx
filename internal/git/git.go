package git

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CloneRepo clones the git repository if it does not already exist
func CloneRepo(repoURL string, path string) error {

	// Check if repo already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Println("Repository already exists. Skipping clone.")
		return nil
	}

	fmt.Println("Cloning repository:", repoURL)

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	fmt.Println("Repository cloned successfully")

	return nil
}

// CommitAndPush stages all changes, commits them, and pushes to remote
func CommitAndPush(repoPath string, message string) error {

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Add all changes
	_, err = w.Add(".")
	if err != nil {
		return err
	}

	// Commit changes
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "sshx",
			Email: "sshx@local",
		},
	})

	if err != nil {
		return err
	}

	// Push changes
	err = repo.Push(&git.PushOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	fmt.Println("Changes pushed to repo")

	return nil
}
