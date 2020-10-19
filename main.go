package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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

	result, err := db.Query("SELECT password, salt FROM user where email=?", values.Get("email"))
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
			var tempPassword string
			var tempSalt string
			result.Scan(&tempPassword, &tempSalt)

			login = checkCredentials(values.Get("password"), tempSalt, tempPassword)
			count = count + 1
		}
	}

	if login {
		fmt.Println("Login successful")
	} else {
		fmt.Println("Login failed")
	}

	fmt.Fprintf(w, "New post was created")
}

func checkCredentials(password string, databaseSalt string, databasePassword string) bool {

	// HMAC-SHA-1 based PBKDF2 key derivation function
	hbts := pbkdf2.Key([]byte(password), []byte(databaseSalt), 10000, 32, sha1.New)
	return databasePassword == base64.StdEncoding.EncodeToString(hbts)

}
