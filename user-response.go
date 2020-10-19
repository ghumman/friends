package main

// UserResponse sent back to user when user endponits are called
type UserResponse struct {
	time     string
	status   int
	err      bool
	message  string
	usersAll []UsersAll
}
