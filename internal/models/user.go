package models

type User struct {
	Name string `json:"name"`
}

func NewUser(name string) User {
	user := User{
		Name: name,
	}

	return user
}
