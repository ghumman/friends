package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/pbkdf2"
)

var db *sql.DB
var err error

// ALPHABET , const used to create salt
const ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const saltLength = 30

func main() {
	db, err = sql.Open("mysql", "ghumman:ghumman@tcp(127.0.0.1:3306)/friends_mysql")
	if err != nil {
		fmt.Printf("err: %v", err)
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/add-user", addUser).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/change-password", changePassword).Methods("POST")
	router.HandleFunc("/forgot-password", forgotPassword).Methods("POST")
	router.HandleFunc("/reset-password", resetPassword).Methods("POST")
	router.HandleFunc("/all-friends", allFriends).Methods("POST")

	router.HandleFunc("/send-message", sendMessage).Methods("POST")
	router.HandleFunc("/messages-user-and-friend", messagesUserAndFriend).Methods("POST")

	http.ListenAndServe(":8000", router)

}

func addUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	values, errParse := url.ParseQuery(string(body))
	if errParse != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := UserResponse{}

	resp.Err = false
	resp.Message = ""

	if values.Get("firstName") == "" {

		resp.Err = true
		resp.Message = "First Name can not be empty. "
	}

	if values.Get("lastName") == "" {

		resp.Err = true
		temp := resp.Message
		resp.Message = temp + "Last Name can not be empty. "
	}

	if values.Get("email") == "" {

		resp.Err = true
		temp := resp.Message
		resp.Message = temp + "Email can not be empty. "
	}

	if values.Get("authType") != "regular" && values.Get("authType") != "special" {

		resp.Err = true
		temp := resp.Message
		resp.Message = temp + "Authentication can only be regular or special. "
	}

	if values.Get("authType") == "regular" && values.Get("password") == "" {
		resp.Err = true
		temp := resp.Message
		resp.Message = temp + "Password can not be empty. "
	}

	if values.Get("authType") == "special" && values.Get("token") == "" {
		resp.Err = true
		temp := resp.Message
		resp.Message = temp + "Token can not be empty. "
	}

	if resp.Err {
		resp.Status = 400
		resp.Time = time.Now().String()
		js, err := json.Marshal(resp)
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(js)
	} else { // The request is valid and required data is present

		// Test if the user already exists
		result, err := db.Query("SELECT password, salt, token FROM user where email=?", values.Get("email"))
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		var count int = 0
		for result.Next() {
			count = count + 1
		}
		if count > 0 { // user already exists
			resp.Message = "User Already Exists"
			resp.Status = 400
			resp.Err = true
			resp.Time = time.Now().String()
			js, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
		} else { // create new user

			// get last ID as unfortuntely Spring Boot didn't create table with id being auto_increment
			lastIDQuery, err := db.Query("SELECT id FROM user ORDER BY id DESC LIMIT 1")
			if err != nil {
				panic(err.Error())
			}
			defer lastIDQuery.Close()
			var lastID int = 0
			for lastIDQuery.Next() {
				lastIDQuery.Scan(&lastID)
			}
			var newID int
			if lastID == 0 {
				newID = 1
			} else {
				newID = lastID + 1
			}

			// create new salt and new hash
			salt := createSalt()
			hash := createHash(values.Get("password"), salt)

			// if creating account using credentials
			if values.Get("authType") == "regular" {

				// now we have newID, salt and hash, and we can creat new user
				createUserQuery, err := db.Query("INSERT INTO `user` (`id`, `auth_type`, `created_at`, `email`,`first_name`, `last_name`, `password`, `salt`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
					newID, 0, time.Now(), values.Get("email"), values.Get("firstName"), values.Get("lastName"), hash, salt)
				if err != nil {
					panic(err.Error())
				}
				defer createUserQuery.Close()

				resp.Message = "User Created"
				resp.Status = 200
				resp.Err = false
				resp.Time = time.Now().String()
				js, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(js)
				sendEmail(values.Get("email"), true, "")

			} else { // else if creating account using OAuth
				// now we have newID, salt and hash, and we can creat new user
				createUserQuery, err := db.Query("INSERT INTO `user` (`id`, `auth_type`, `created_at`, `email`,`first_name`, `last_name`, `token`) VALUES (?, ?, ?, ?, ?, ?, ?)",
					newID, 1, time.Now(), values.Get("email"), values.Get("firstName"), values.Get("lastName"), values.Get("token"))
				if err != nil {
					panic(err.Error())
				}
				defer createUserQuery.Close()

				resp.Message = "User Created"
				resp.Status = 200
				resp.Err = false
				resp.Time = time.Now().String()
				js, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(js)
				sendEmail(values.Get("email"), true, "")
			} // else if creating account using OAuth ends
		} // create new user ends
	} // The request is valid and required data is present ends
}

func login(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	values, errParse := url.ParseQuery(string(body))
	if errParse != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := db.Query("SELECT password, salt, token FROM user where email=?", values.Get("email"))
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var count int = 0

	var login bool = false
	var loginCredentialsFailed bool = false
	var tempPassword, tempSalt, tempToken sql.NullString

	for result.Next() {

		if count > 1 {
			login = false
			break
		} else {

			result.Scan(&tempPassword, &tempSalt, &tempToken)

			if values.Get("authType") == "regular" {

				login = checkCredentials(values.Get("password"), tempSalt.String, tempPassword.String)
				if login == false {
					loginCredentialsFailed = true
				}
				count = count + 1
			} else if values.Get("authType") == "special" {
				// account is setup using OAuth
				if values.Get("token") == tempToken.String {
					login = true
				} else {
					login = false
					loginCredentialsFailed = true
				}
			}
		}
	}

	resp := UserResponse{}

	if login {
		resp.Message = "Logged In"
		resp.Err = false
		resp.Status = 200
		resp.Time = time.Now().String()
	} else {

		if loginCredentialsFailed {
			resp.Message = "Login Failed"
		} else {
			resp.Message = "User Does Not Exist"
		}

		resp.Err = true
		resp.Status = 400
		resp.Time = time.Now().String()

	}
	js, err := json.Marshal(resp)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func checkCredentials(password string, databaseSalt string, databasePassword string) bool {

	// HMAC-SHA-1 based PBKDF2 key derivation function
	hbts := pbkdf2.Key([]byte(password), []byte(databaseSalt), 10000, 32, sha1.New)
	return databasePassword == base64.StdEncoding.EncodeToString(hbts)

}

func createSalt() string {

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	var ret string
	for i := 0; i < saltLength; i++ {
		str := ALPHABET[seededRand.Intn(len(ALPHABET)-1)]
		ret += string(str)
	}

	return ret

}

func createHash(password string, databaseSalt string) string {
	hbts := pbkdf2.Key([]byte(password), []byte(databaseSalt), 10000, 32, sha1.New)
	return base64.StdEncoding.EncodeToString(hbts)
}

func sendEmail(toEmail string, accountCreated bool, token string) {

	// Choose auth method and set it up
	auth := smtp.PlainAuth("", "your-email.com", "your-email-password", "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{toEmail}

	if accountCreated {

		msg := []byte("To: " + toEmail + "\r\n" +
			"Subject: Welcome to Friends\r\n" +
			"\r\n" +
			"New Account Created\r\n")
		err := smtp.SendMail("smtp.gmail.com:587", auth, "your-email.com", to, msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := []byte("To: " + toEmail + "\r\n" +
			"Subject: Password Reset Request\r\n" +
			"\r\n" +
			"To reset your password, click the link below:\n" +
			"http://localhost:3000/#/reset-password?token=" + token +
			"\r\n")
		err := smtp.SendMail("smtp.gmail.com:587", auth, "your-email.com", to, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func changePassword(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	values, errParse := url.ParseQuery(string(body))
	if errParse != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := UserResponse{}
	var tempPassword, tempSalt, tempToken sql.NullString
	var tempID sql.NullInt64

	if values.Get("authType") == "special" {
		resp.Message = "Logged in using OAuth"
		resp.Err = true
		resp.Status = 400
		resp.Time = time.Now().String()

		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	} else if values.Get("authType") == "regular" {
		result, err := db.Query("SELECT id, password, salt, token FROM user where email=?", values.Get("email"))
		if err != nil {
			panic(err.Error())
		}

		defer result.Close()

		var count int = 0

		var login bool = false

		for result.Next() {

			if count > 1 {
				login = false
				break
			} else {

				result.Scan(&tempID, &tempPassword, &tempSalt, &tempToken)

				login = checkCredentials(values.Get("password"), tempSalt.String, tempPassword.String)
				count = count + 1

			}
		}

		if count == 0 {
			resp.Message = "User Does Not Exist"
			resp.Err = true
			resp.Status = 400
		} else if login {

			// create new salt and new hash
			salt := createSalt()
			hash := createHash(values.Get("newPassword"), salt)

			updateUserQuery, err := db.Query("UPDATE `user` SET salt=?, password=? where id=?", salt, hash, tempID)
			if err != nil {
				panic(err.Error())
			}
			defer updateUserQuery.Close()

			resp.Message = "Password changed"
			resp.Err = false
			resp.Status = 200
		} else {
			resp.Message = "Original password not right"
			resp.Err = true
			resp.Status = 400
		}

		resp.Time = time.Now().String()
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(js)
	} else {
		resp.Message = "Authentication Type not right"
		resp.Err = true
		resp.Status = 400
		resp.Time = time.Now().String()

		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	}
}

func forgotPassword(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	values, errParse := url.ParseQuery(string(body))
	if errParse != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := UserResponse{}
	var tempAuthType, tempID int64

	if values.Get("email") == "" {
		resp.Message = "Email is required"
		resp.Err = true
		resp.Status = 400
		resp.Time = time.Now().String()

		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	} else {
		result, err := db.Query("SELECT id, auth_type FROM user where email=?", values.Get("email"))
		if err != nil {
			panic(err.Error())
		}

		defer result.Close()

		var count int = 0

		for result.Next() {

			result.Scan(&tempID, &tempAuthType)
			count = count + 1
		}

		if count == 0 || tempAuthType == 1 {
			resp.Message = "No account or logged in using social media"
			resp.Err = true
			resp.Status = 400
		} else if count == 1 {

			uuid := uuid.New()

			updateUserQuery, err := db.Query("UPDATE `user` SET reset_token=? where id=?", uuid.String(), tempID)
			if err != nil {
				panic(err.Error())
			}
			defer updateUserQuery.Close()

			// send email
			sendEmail(values.Get("email"), false, uuid.String())
			resp.Message = "Reset password is sent"
			resp.Err = false
			resp.Status = 200
		}

		resp.Time = time.Now().String()
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(js)
	}
}

func resetPassword(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	values, errParse := url.ParseQuery(string(body))
	if errParse != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := UserResponse{}

	var tempID sql.NullInt64

	// Test if the user already exists
	result, err := db.Query("SELECT id FROM user where reset_token=?", values.Get("token"))
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var count int = 0
	for result.Next() {
		result.Scan(&tempID)
		count = count + 1
	}
	if count == 0 { // user already exists
		resp.Message = "Token is not valid"
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
	} else if count == 1 { // update password

		salt := createSalt()
		hash := createHash(values.Get("password"), salt)

		updateUserQuery, err := db.Query("UPDATE `user` SET salt=?, password=?, reset_token=? where id=?", salt, hash, nil, tempID)
		if err != nil {
			panic(err.Error())
		}
		defer updateUserQuery.Close()

		resp.Message = "Password successfully reset"
		resp.Status = 200
		resp.Err = false
		resp.Time = time.Now().String()
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	} else {
		resp.Message = "Multiple users found against this reset token"
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

func allFriends(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	values, errParse := url.ParseQuery(string(body))
	if errParse != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := db.Query("SELECT password, salt, token FROM user where email=?", values.Get("email"))
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var count int = 0

	var login bool = false
	var loginCredentialsFailed bool = false
	var tempPassword, tempSalt, tempToken sql.NullString

	resp := UserResponse{}

	for result.Next() {

		if count > 1 {
			login = false
			break
		} else {

			result.Scan(&tempPassword, &tempSalt, &tempToken)

			if values.Get("authType") == "regular" {

				login = checkCredentials(values.Get("password"), tempSalt.String, tempPassword.String)
				if login == false {
					loginCredentialsFailed = true
				} else {
					resultFriends, err := db.Query("SELECT first_name, last_name, email FROM user where email!=?", values.Get("email"))
					if err != nil {
						panic(err.Error())
					}

					defer resultFriends.Close()

					var tempFirstName, tempLastName, tempEmail sql.NullString

					for resultFriends.Next() {

						resultFriends.Scan(&tempFirstName, &tempLastName, &tempEmail)

						userData := UsersAll{
							tempFirstName.String,
							tempLastName.String,
							tempEmail.String,
						}
						resp.Users = append(resp.Users, userData)
					}
				}
				count = count + 1
			} else if values.Get("authType") == "special" {
				// account is setup using OAuth
				if values.Get("token") == tempToken.String {
					login = true

				} else {
					login = false
					loginCredentialsFailed = true
				}
			}
		}
	}

	if login {
		resp.Message = "Friends attached"
		resp.Err = false
		resp.Status = 200
		resp.Time = time.Now().String()
	} else {

		if loginCredentialsFailed {
			resp.Message = "Login Failed"
		} else {
			resp.Message = "User Does Not Exist"
		}

		resp.Err = true
		resp.Status = 400
		resp.Time = time.Now().String()

	}
	js, err := json.Marshal(resp)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

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

	resultSenderError := db.QueryRow("SELECT id, password, salt, token, auth_type FROM user where email=?", values.Get("messageFromEmail")).Scan(
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

	resultReceiverError := db.QueryRow("SELECT id, password, salt, token, auth_type FROM user where email=?", values.Get("messageToEmail")).Scan(
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
			createMessageQuery, err := db.Query("INSERT INTO `message` (`id`, `message`, `sent_at`, `message_from_id`,`message_to_id`) VALUES (?, ?, ?, ?, ?)",
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
			createMessageQuery, err := db.Query("INSERT INTO `message` (`id`, `message`, `sent_at`, `message_from_id`,`message_to_id`) VALUES (?, ?, ?, ?, ?)",
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

	resultSenderError := db.QueryRow("SELECT id, password, salt, token, auth_type FROM user where email=?", values.Get("userEmail")).Scan(
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

	resultReceiverError := db.QueryRow("SELECT id, password, salt, token, auth_type FROM user where email=?", values.Get("friendEmail")).Scan(
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
				"SELECT message, message_from_id, message_to_id, sent_at  FROM message m WHERE (m.message_from_id = ? and m.message_to_id = ?) or (m.message_from_id = ? and m.message_to_id = ?) order by m.sent_at",
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
				"SELECT message, message_from_id, message_to_id, sent_at  FROM message m WHERE (m.message_from_id = ? and m.message_to_id = ?) or (m.message_from_id = ? and m.message_to_id = ?) order by m.sent_at",
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
