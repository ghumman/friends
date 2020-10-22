package main

// MessageResponse ...
type MessageResponse struct {
	Time    string        `json:"time"`
	Status  int           `json:"status"`
	Err     bool          `json:"err"`
	Message string        `json:"message"`
	Msgs    []MessagesAll `json:"msgs"`
}
