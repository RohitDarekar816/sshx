package git

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CommitAndPush(repoPath string, message string) error {

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	w.Add(".")

	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "sshx",
			Email: "sshx@local",
		},
	})

	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{})

	if err != nil {
		return err
	}

	fmt.Println("Changes pushed to repo")

	return nil
}
