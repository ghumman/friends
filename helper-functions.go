package main

import (
	"crypto/sha1"
	"encoding/base64"
	"log"
	"math/rand"
	"net/smtp"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

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
