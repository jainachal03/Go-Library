package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type authHandler struct {
	l      *log.Logger
	client *mongo.Client
}

func NewAuthHandler(l *log.Logger, client *mongo.Client) *authHandler {
	return &authHandler{l, client}
}

type User struct {
	Email    string `bson:"Email"`
	Password string `bson:"Password"`
}

// signup handler
func (a *authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t User
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	// before inserting make sure that there is already no user as such?
	client := a.client
	coll := client.Database("Books").Collection("User")
	Email := t.Email
	Password := t.Password
	a.l.Println(Email, Password)
	doc := bson.D{{"Email", Email}, {"Password", Password}}

	var result bson.M
	err = coll.FindOne(context.TODO(), doc).Decode(&result)

	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(" You are already a registered User!"))
		return
	}
	// before inserting check if the use
	res, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(" You are already a registered User!"))
		return
	}
	a.l.Println(res.InsertedID)
	w.Write([]byte(" Sign up Successful"))
	return

}

// login handler
func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var t User
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	email := t.Email
	password := t.Password
	a.l.Println(email)
	a.l.Println(password)

	if len(email) > 0 {
		if len(password) > 0 {
			coll := a.client.Database("Books").Collection("Auth")
			var result bson.D
			err := coll.FindOne(context.TODO(), bson.D{{"Email", email}, {"Password", password}}).Decode(&result)
			if err == mongo.ErrNoDocuments {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - Something bad happened!"))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Successful Sign in"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {

// 	email := r.URL.Query().Get("email")
// 	password := r.URL.Query().Get("password")

// 	a.l.Println(email, password)
// 	coll := a.client.Database("bookHandler").Collection("Auth")
// 	var result bson.D
// 	err := coll.FindOne(context.TODO(), bson.D{{"email", email}, {"password", password}}).Decode(&result)
// 	if err == mongo.ErrNoDocuments {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("400 - Something bad happened!"))
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	// There was one so sucessful

// 	// now check if this exists or not.
// }
