package git

import (
	"fmt"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CloneRepo clones the git repository if it does not already exist
func CloneRepo(repoURL, path string) error {

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

// CommitAndPush stages all changes, commits them, and pushes to remote.
// The commit is attributed to the given author so the shared repo carries a
// meaningful per-user audit trail; empty values fall back to a local identity.
func CommitAndPush(repoPath, message, authorName, authorEmail string) error {

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

	if authorName == "" {
		authorName = "sshx"
	}
	if authorEmail == "" {
		authorEmail = "sshx@local"
	}
	sig := &object.Signature{
		Name:  authorName,
		Email: authorEmail,
		When:  time.Now(),
	}

	// Commit changes
	_, err = w.Commit(message, &git.CommitOptions{
		Author:    sig,
		Committer: sig,
	})

	if err != nil {
		return err
	}

	// Push changes
	err = repo.Push(&git.PushOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		if isNonFastForward(err) {
			if pullErr := PullRebase(repo); pullErr != nil {
				return fmt.Errorf("push rejected and rebase failed: %w", pullErr)
			}

			err = repo.Push(&git.PushOptions{})
		}
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
	}

	fmt.Println("Changes pushed to repo")

	return nil
}

func isNonFastForward(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "non-fast-forward")
}
