package user

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Joined   string `json:"joined"`
}

type UsersFile struct {
	Users []User `json:"users"`
}
