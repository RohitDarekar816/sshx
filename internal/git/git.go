package git

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
)

func CloneRepo(repoURL string, path string) error {

	fmt.Println("Cloning repository...")

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: repoURL,
	})

	if err != nil {
		return err
	}

	fmt.Println("Repository cloned successfully")

	return nil
}
