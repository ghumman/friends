package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func sendMessage(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	values, errParse := url.ParseQuery(string(body))
	if errParse != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tempSenderPassword, tempReceiverPassword, tempSenderSalt, tempReceiverSalt, tempSenderToken, tempReceiverToken sql.NullString
	var tempSenderID, tempReceiverID, tempSenderAuthType, tempReceiverAuthType sql.NullInt64
	resp := MessageResponse{}

	resultSenderError := db.QueryRow("SELECT id, password, salt, token, auth_type FROM users where email=$1", values.Get("messageFromEmail")).Scan(
		&tempSenderID, &tempSenderPassword, &tempSenderSalt, &tempSenderToken, &tempSenderAuthType,
	)
	switch {
	case resultSenderError == sql.ErrNoRows:
		resp.Message = "Sender Does Not Exist"
		resp.Status = 400
		resp.Err = true
		resp.Time = time.Now().String()
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	case resultSenderError != nil:
		log.Fatal(err)
	}

	resultReceiverError := db.QueryRow("SELECT id, password, salt, token, auth_type FROM users where email=$1", values.Get("messageToEmail")).Scan(
		&tempReceiverID, &tempReceiverPassword, &tempReceiverSalt, &tempReceiverToken, &tempReceiverAuthType,
	)
	switch {
	case resultReceiverError == sql.ErrNoRows:
		resp.Message = "Receiver Does Not Exist"
		resp.Status = 400
		resp.Err = true
		resp.Time = time.Now().String()
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	case resultReceiverError != nil:
		log.Fatal(err)
	}

	var tempMessageID int64
	lastMessageError := db.QueryRow("SELECT id FROM message ORDER BY id DESC LIMIT 1").Scan(&tempMessageID)

	switch {
	case lastMessageError == sql.ErrNoRows:
		tempMessageID = 1
	default:
		tempMessageID = tempMessageID + 1
	}

	if tempSenderAuthType.Int64 == 0 {

		login := checkCredentials(values.Get("password"), tempSenderSalt.String, tempSenderPassword.String)
		if login {
			createMessageQuery, err := db.Query("INSERT INTO message (id, message, sent_at, message_from_id, message_to_id) VALUES ($1, $2, $3, $4, $5)",
				tempMessageID, values.Get("message"), time.Now(), tempSenderID, tempReceiverID)
			if err != nil {
				panic(err.Error())
			}
			defer createMessageQuery.Close()

			resp.Message = "Message sent"
			resp.Status = 200
			resp.Err = false
			resp.Time = time.Now().String()
			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
			return
		} else {
			resp.Message = "Login Failed"
			resp.Status = 400
			resp.Err = true
			resp.Time = time.Now().String()
			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
			return
		}

	} else if tempSenderAuthType.Int64 == 1 {
		if values.Get("token") == tempSenderToken.String {
			createMessageQuery, err := db.Query("INSERT INTO message (id, message, sent_at, message_from_id,message_to_id) VALUES ($1, $2, $3, $4, $5)",
				tempMessageID, values.Get("message"), time.Now(), tempSenderID, tempReceiverID)
			if err != nil {
				panic(err.Error())
			}
			defer createMessageQuery.Close()

			resp.Message = "Message sent"
			resp.Status = 200
			resp.Err = false
			resp.Time = time.Now().String()
			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
			return
		} else {
			resp.Message = "Login Failed"
			resp.Status = 400
			resp.Err = true
			resp.Time = time.Now().String()
			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
			return
		}
	}
}

func messagesUserAndFriend(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	values, errParse := url.ParseQuery(string(body))
	if errParse != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tempSenderPassword, tempReceiverPassword, tempSenderSalt, tempReceiverSalt, tempSenderToken, tempReceiverToken sql.NullString
	var tempSenderID, tempReceiverID, tempSenderAuthType, tempReceiverAuthType sql.NullInt64
	resp := MessageResponse{}

	resultSenderError := db.QueryRow("SELECT id, password, salt, token, auth_type FROM users where email=$1", values.Get("userEmail")).Scan(
		&tempSenderID, &tempSenderPassword, &tempSenderSalt, &tempSenderToken, &tempSenderAuthType,
	)
	switch {
	case resultSenderError == sql.ErrNoRows:
		resp.Message = "Sender Does Not Exist"
		resp.Status = 400
		resp.Err = true
		resp.Time = time.Now().String()
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	case resultSenderError != nil:
		log.Fatal(err)
	}

	resultReceiverError := db.QueryRow("SELECT id, password, salt, token, auth_type FROM users where email=$1", values.Get("friendEmail")).Scan(
		&tempReceiverID, &tempReceiverPassword, &tempReceiverSalt, &tempReceiverToken, &tempReceiverAuthType,
	)
	switch {
	case resultReceiverError == sql.ErrNoRows:
		resp.Message = "Receiver Does Not Exist"
		resp.Status = 400
		resp.Err = true
		resp.Time = time.Now().String()
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	case resultReceiverError != nil:
		log.Fatal(err)
	}
	if tempSenderAuthType.Int64 == 0 {

		login := checkCredentials(values.Get("password"), tempSenderSalt.String, tempSenderPassword.String)
		if login {

			resultMessages, err := db.Query(
				"SELECT message, message_from_id, message_to_id, sent_at  FROM message m WHERE (m.message_from_id = $1 and m.message_to_id = $2) or (m.message_from_id = $3 and m.message_to_id = $4) order by m.sent_at",
				tempSenderID, tempReceiverID, tempReceiverID, tempSenderID)
			if err != nil {
				panic(err.Error())
			}

			defer resultMessages.Close()

			var tempMessage, tempMessageFromEmail, tempMessageToEmail, tempSentAt sql.NullString
			var tempMessageFromEmailID, tempMessageToEmailID sql.NullInt64
			for resultMessages.Next() {

				resultMessages.Scan(&tempMessage, &tempMessageFromEmailID, &tempMessageToEmailID, &tempSentAt)

				if (tempMessageFromEmailID == tempSenderID) && (tempMessageToEmailID == tempReceiverID) {
					tempMessageFromEmail = sql.NullString{String: values.Get("userEmail"), Valid: true}
					tempMessageToEmail = sql.NullString{String: values.Get("friendEmail"), Valid: true}
				} else if (tempMessageFromEmailID == tempReceiverID) && (tempMessageToEmailID == tempSenderID) {
					tempMessageFromEmail = sql.NullString{String: values.Get("friendEmail"), Valid: true}
					tempMessageToEmail = sql.NullString{String: values.Get("userEmail"), Valid: true}
				}
				messageData := MessagesAll{
					tempMessage.String,
					tempMessageFromEmail.String,
					tempMessageToEmail.String,
					tempSentAt.String,
				}
				resp.Msgs = append(resp.Msgs, messageData)
			}

			resp.Message = "Messages attached"
			resp.Status = 200
			resp.Err = false
			resp.Time = time.Now().String()
			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
			return
		} else {
			resp.Message = "Login Failed"
			resp.Status = 400
			resp.Err = true
			resp.Time = time.Now().String()
			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
			return
		}

	} else if tempSenderAuthType.Int64 == 1 {
		if values.Get("token") == tempSenderToken.String {

			resultMessages, err := db.Query(
				"SELECT message, message_from_id, message_to_id, sent_at  FROM message m WHERE (m.message_from_id = $1 and m.message_to_id = $2) or (m.message_from_id = $3 and m.message_to_id = $4) order by m.sent_at",
				tempSenderID, tempReceiverID, tempReceiverID, tempSenderID)
			if err != nil {
				panic(err.Error())
			}

			defer resultMessages.Close()

			var tempMessage, tempMessageFromEmail, tempMessageToEmail, tempSentAt sql.NullString
			var tempMessageFromEmailID, tempMessageToEmailID sql.NullInt64
			for resultMessages.Next() {

				resultMessages.Scan(&tempMessage, &tempMessageFromEmailID, &tempMessageToEmailID, &tempSentAt)

				if (tempMessageFromEmailID == tempSenderID) && (tempMessageToEmailID == tempReceiverID) {
					tempMessageFromEmail = sql.NullString{String: values.Get("userEmail"), Valid: true}
					tempMessageToEmail = sql.NullString{String: values.Get("friendEmail"), Valid: true}
				} else if (tempMessageFromEmailID == tempReceiverID) && (tempMessageToEmailID == tempSenderID) {
					tempMessageFromEmail = sql.NullString{String: values.Get("friendEmail"), Valid: true}
					tempMessageToEmail = sql.NullString{String: values.Get("userEmail"), Valid: true}
				}
				messageData := MessagesAll{
					tempMessage.String,
					tempMessageFromEmail.String,
					tempMessageToEmail.String,
					tempSentAt.String,
				}
				resp.Msgs = append(resp.Msgs, messageData)
			}

			resp.Message = "Message sent"
			resp.Status = 200
			resp.Err = false
			resp.Time = time.Now().String()
			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
			return
		} else {
			resp.Message = "Login Failed"
			resp.Status = 400
			resp.Err = true
			resp.Time = time.Now().String()
			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
		}
	}
}
