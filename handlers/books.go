package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jainachal03/lib/data"
	"go.mongodb.org/mongo-driver/mongo"
)

type bookHandler struct {
	l      *log.Logger
	client *mongo.Client
}

func NewBooks(l *log.Logger, client *mongo.Client) *bookHandler {
	return &bookHandler{l, client}
}
func (b *bookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	b.l.Println("Handle GET Products")
	data := data.GetBooks(b.client)
	t, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Cotent-Type", "application/json")
	w.Write(t)
}
func (b *bookHandler) GetBooksByAuthor(w http.ResponseWriter, r *http.Request) {
	b.l.Println("Handle GET Books by Author")

	author := mux.Vars(r)["Author"]
	b.l.Println(" the author is ", author)
	res, err := data.GetBooksByAuthor(b.client, author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b.l.Println(res)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(" The books have been fetched"))

}

// called by the handler.

// very simple function to quer validation.
func bookValidate(b *data.Book) bool {
	if len(b.Author) == 0 {
		return false
	}
	if len(b.Title) == 0 {
		return false
	}
	if b.CurrentCopies <= 0 {
		return false
	}
	if b.TotalCopies <= 0 {
		return false
	}
	return true
}

func (b *bookHandler) AddBooks(w http.ResponseWriter, r *http.Request) {
	b.l.Println("Handle Post Products")
	// extract the result from this post request.

	decoder := json.NewDecoder(r.Body)
	var t data.Book
	err := decoder.Decode(&t)
	// make sure if any ting is less than bad.
	if err != nil {
		panic(err)
	}
	// validate here.
	res := bookValidate(&t)

	if res == false {
		if err != nil {
			panic(err)
		}
	}

	err = data.AddBooks(b.client, &t)
	if err != nil {
		b.l.Println("Bad Request Loool")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

	// first extract the details ?
}
