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
	UpdateUser(id int, user User) error
	DeleteUser(id int) error
}

type InMemoryStorage struct {
	users map[int]User
}

type UserService struct {
	userStorage UserStorage
}

func (u *UserService) CreateUser(user User) error {
	return u.userStorage.CreateUser(user)
}

func (u *UserService) GetUser(id int) (User, error) {
	return u.userStorage.GetUser(id)
}

func(u *UserService) UpdateUser(id int, user User) error {
	return u.userStorage.UpdateUser(id, user)
}

func(u *UserService) DeleteUser(id int) error {
	return u.userStorage.DeleteUser(id)
}

func (s *InMemoryStorage) CreateUser(user User) error {
	s.users[user.id] = user
	return nil
}

func (s *InMemoryStorage) GetUser(id int) (User, error) {
	user, exists := s.users[id]
	if !exists {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *InMemoryStorage) UpdateUser(id int, user User) error {
	if user.id != id {
		return fmt.Errorf("user id mismatch")
	}
	_, exists := s.users[id]
	if (!exists) {
		return fmt.Errorf("user does not exist")
	}
	s.users[id] = user
	return nil
}

func (s *InMemoryStorage) DeleteUser(id int) error {
	_, exists := s.users[id]
	if (!exists) {
		return fmt.Errorf("user does not exists")
	}
	delete(s.users, id)
	return nil
}

func main() {
	s := &InMemoryStorage{
		users : make(map[int]User),
	}
	service := UserService{userStorage: s}

	userInput := User{id: 1, name: "John", email: "john@example.com"}
	service.CreateUser(userInput)
	userOutput, err := service.GetUser(userInput.id)
	if err != nil {
		fmt.Println("Error getting users:", err)
		return
	}
	fmt.Println("user received back: " , userOutput.name)

	newUserInput:= User{id: 1, name: "Smith", email: "smith@example.com"}

	service.UpdateUser(1, newUserInput);

	service.DeleteUser(1);
}