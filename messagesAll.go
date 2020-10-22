package main

// MessagesAll ...
type MessagesAll struct {
	Message          string `json:"message"`
	MessageFromEmail string `json:"messageFromEmail"`
	MessageToEmail   string `json:"messageToEmail"`
	SentAt           string `json:"sentAt"`
}
