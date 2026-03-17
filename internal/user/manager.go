package user

import (
	"encoding/json"
	"os"
	"time"
)

func LoadUsers(path string) (*UsersFile, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &UsersFile{Users: []User{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var users UsersFile
	err = json.Unmarshal(data, &users)

	return &users, err
}

func SaveUsers(path string, users *UsersFile) error {

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func UserExists(users *UsersFile, email string) bool {

	for _, u := range users.Users {
		if u.Email == email {
			return true
		}
	}

	return false
}

func AddUser(users *UsersFile, username, email string) {

	users.Users = append(users.Users, User{
		Username: username,
		Email:    email,
		Joined:   time.Now().Format("2006-01-02"),
	})
}
