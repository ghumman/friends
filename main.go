package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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
