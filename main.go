package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection, messageCollection *mongo.Collection
var db *sql.DB
var err error

// ALPHABET , const used to create salt
const ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const saltLength = 30

func main() {
	// db, err = sql.Open("mysql", "ghumman:ghumman@tcp(127.0.0.1:3306)/friends_mysql")
	// if err != nil {
	// 	fmt.Printf("err: %v", err)
	// 	panic(err.Error())
	// }
	// defer db.Close()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("friends_mongo")

	userCollection = db.Collection("user")
	messageCollection = db.Collection("message")

	_ = userCollection
	_ = messageCollection

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
