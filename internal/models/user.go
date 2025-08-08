package models

import (
	"fmt"
	"sync"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUser(name string, email string) User {
	user := User{
		Name:  name,
		Email: email,
	}

	return user
}

type UserStore struct {
	users map[string]*User
	mu    sync.RWMutex
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*User),
	}
}

func (us *UserStore) AddUser(user *User) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.users[user.Email]; exists {
		return fmt.Errorf("User with email %s already exists", user.Email)
	}

	us.users[user.Email] = user
	return nil
}

func (us *UserStore) RemoveUser(user *User) {
	us.mu.Lock()
	defer us.mu.Unlock()

	delete(us.users, user.Email)
}
