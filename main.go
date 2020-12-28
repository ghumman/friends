package main

import (
	"database/sql"
	"fmt"
	"net/http"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

// ALPHABET , const used to create salt
const ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const saltLength = 30

const (
	host     = "localhost"
	user     = "postgres"
	port     = 5432
	password = "postgres"
	dbname   = "friends_psql"
)

func main() {

	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s",
	// 	host, port, user, password, dbname)

	// // db, err = sql.Open("mysql", "ghumman:ghumman@tcp(127.0.0.1:3306)/friends_mysql")

	// fmt.Println("value of psqlInfo: ", psqlInfo)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	fmt.Printf("err: %v", err)
	// 	panic(err.Error())
	// }
	// defer db.Close()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Successfully connected!")

	// userSql := "SELECT token FROM users WHERE email = $1"

	// var myPassword sql.NullString
	// err = db.QueryRow(userSql, "test001@gmail.com").Scan(&myPassword)
	// if err != nil {
	// 	log.Fatal("Failed to execute query: ", err)
	// }

	// fmt.Printf("Hi %s, welcome back!\n", myPassword)

	// userSql2 := "SELECT password, salt, token FROM users WHERE email = $1"

	// result, _ := db.Query(userSql2, "test001@gmail.com")

	// fmt.Println("value of result: ", result)
	// // if err != nil {
	// // 	log.Fatal("Failed to execute query: ", err)
	// // }

	// fmt.Println("After query")

	// fmt.Println("value of db: ", db)

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
