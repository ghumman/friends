package main

import (
	"fmt"
)

type User struct {
	id int
	name string
	email string
}

type UserStorage interface {
	CreateUser(user User) error
	GetUser(id int) (User, error)
}

type InMemoryStorage struct {
	users map[int]User
}

type UserService struct {
	userStorage UserStorage
}

func (u UserService) CreateUser(user User) error {
	return u.userStorage.CreateUser(user)
}

func (u UserService) GetUser(id int) (User, error) {
	return u.userStorage.GetUser(id)
}

func (s *InMemoryStorage) CreateUser(user User) error {
	s.users[user.id] = user
	return nil
}

func (s *InMemoryStorage) GetUser(id int) (User, error) {
	user, exists := s.users[id]
	if !exists {
		return User{}, fmt.Errorf("User not found")
	}
	return user, nil
}

func main() {
	s := &InMemoryStorage{
		users : make(map[int]User),
	}
	service := UserService{userStorage: s}

	userInput := User{id: 1, name: "John", email: "john@example.com"}
	service.CreateUser(userInput)
	userOutput, _ := service.GetUser(userInput.id)
	fmt.Println("user received back: " , userOutput.name)
}