package git

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
)

func PullRebase(repo *git.Repository) error {

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Rebase:     true,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("pull failed: %w", err)
	}

	return nil
}
