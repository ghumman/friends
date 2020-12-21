package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

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

		var result User
		err = userCollection.FindOne(context.TODO(), bson.D{{"email", values.Get("email")}}).Decode(&result)

		if result.Email != "" {
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

			// create new salt and new hash
			salt := createSalt()
			hash := createHash(values.Get("password"), salt)

			// if creating account using credentials
			if values.Get("authType") == "regular" {

				// now we have newID, salt and hash, and we can creat new user

				_, err := userCollection.InsertOne(context.TODO(), bson.M{"firstName": values.Get("firstName"), "lastName": values.Get("lastName"), "createdAt": time.Now(), "salt": salt, "password": hash, "email": values.Get("email"), "authType": 0})
				if err != nil {
					log.Fatal(err)
				}

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
				// sendEmail(values.Get("email"), true, "")

			} else { // else if creating account using OAuth
				// now we have newID, salt and hash, and we can creat new user

				_, err := userCollection.InsertOne(context.TODO(), bson.M{"firstName": values.Get("firstName"), "lastName": values.Get("lastName"), "createdAt": time.Now(), "salt": salt, "password": hash, "email": values.Get("email"), "authType": 0})

				if err != nil {
					log.Fatal(err)
				}

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
				// sendEmail(values.Get("email"), true, "")
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

	var result User
	err = userCollection.FindOne(context.TODO(), bson.D{{"email", values.Get("email")}}).Decode(&result)

	var login bool = false
	var loginCredentialsFailed bool = false

	if values.Get("authType") == "regular" {

		login = checkCredentials(values.Get("password"), result.Salt, result.Password)
		if login == false {
			loginCredentialsFailed = true
		}
	} else if values.Get("authType") == "special" {
		// account is setup using OAuth
		if values.Get("token") == result.Token {
			login = true
		} else {
			login = false
			loginCredentialsFailed = true
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
