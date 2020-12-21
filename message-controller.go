package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	resp := MessageResponse{}
	var resultSender, resultReceiver User

	// check if sender exists
	userCollection.FindOne(context.TODO(), bson.D{{"email", values.Get("messageFromEmail")}}).Decode(&resultSender)
	if resultSender.Email == "" {
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
	} else {
		// check if receiver exists
		userCollection.FindOne(context.TODO(), bson.D{{"email", values.Get("messageToEmail")}}).Decode(&resultReceiver)
		if resultReceiver.Email == "" {
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
		} else {
			// check if user can login
			login := checkCredentials(values.Get("password"), resultSender.Salt, resultSender.Password)
			if !login {
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
			} else {
				// everything looks good, send the message
				_, err := messageCollection.InsertOne(context.TODO(), bson.M{"message": values.Get("message"), "messageFrom": resultSender, "messageTo": resultReceiver, "sentAt": time.Now()})
				if err != nil {
					log.Fatal(err)
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
			}
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

	resp := MessageResponse{}
	var resultSender, resultReceiver User

	// check if sender exists
	userCollection.FindOne(context.TODO(), bson.D{{"email", values.Get("userEmail")}}).Decode(&resultSender)
	if resultSender.Email == "" {
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
	} else {
		// check if receiver exists
		userCollection.FindOne(context.TODO(), bson.D{{"email", values.Get("friendEmail")}}).Decode(&resultReceiver)
		if resultReceiver.Email == "" {
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
		} else {
			// check if user can login
			login := checkCredentials(values.Get("password"), resultSender.Salt, resultSender.Password)
			if !login {
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
			} else {
				// everything looks good, send back all the messages

				cur, err := messageCollection.Find(context.Background(), bson.M{"$or": []interface{}{bson.M{"$and": []interface{}{bson.M{"messageFrom": resultSender}, bson.M{"messageTo": resultReceiver}}}, bson.M{"$and": []interface{}{bson.M{"messageFrom": resultReceiver}, bson.M{"messageTo": resultSender}}}}})
				if err != nil {
					log.Fatal(err)
				}
				defer cur.Close(context.Background())
				for cur.Next(context.Background()) {
					var currentResult Message
					err := cur.Decode(&currentResult)
					if err != nil {
						log.Fatal(err)
					}

					messageData := MessagesAll{
						currentResult.Message,
						currentResult.MessageFrom.Email,
						currentResult.MessageTo.Email,
						currentResult.SentAt.String(),
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
			}
		}
	}
}
