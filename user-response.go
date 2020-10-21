package main

// UserResponse ...
type UserResponse struct {
	Time    string     `json:"time"`
	Status  int        `json:"status"`
	Err     bool       `json:"err"`
	Message string     `json:"message"`
	Users   []UsersAll `json:"usersAll"`
}
