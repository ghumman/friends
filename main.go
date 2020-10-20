package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/pbkdf2"
)

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "ghumman:ghumman@tcp(127.0.0.1:3306)/friends_mysql")
	if err != nil {
		fmt.Printf("err: %v", err)
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/login", login).Methods("POST")

	http.ListenAndServe(":8000", router)

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

	for result.Next() {

		if count > 1 {
			login = false
			break
		} else {
			var tempPassword, tempSalt, tempToken string

			if values.Get("authType") == "regular" {

				result.Scan(&tempPassword, &tempSalt, &tempToken)

				login = checkCredentials(values.Get("password"), tempSalt, tempPassword)
				if login == false {
					loginCredentialsFailed = true
				}
				count = count + 1
			} else if values.Get("authType") == "special" {
				// account is setup using OAuth
				if values.Get("token") == tempToken {
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
