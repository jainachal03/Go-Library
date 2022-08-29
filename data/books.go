package data

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Book struct {
	Title         string `json:"Title"`
	Author        string `json:"Author"`
	Genre         string `json:"Genre"`
	TotalCopies   int32  `json:"Totalcopies"`
	CurrentCopies int32  `json:"Currentcopies"`
}

func (b *Book) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(b)
}

func GetBooks(client *mongo.Client) []bson.D {
	coll := client.Database("Books").Collection("Cluster0")
	var results []bson.D // result of type

	cur, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem bson.D
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)

	}
	return results
}

// can be many books :-)
func getBookFrom(d primitive.D) Book {
	var author, title, genre string
	var totalcopies, currentcopies int32
	for _, val := range d {
		if val.Key == "Author" {
			author = val.Value.(string)
		}
		if val.Key == "Title" {
			title = val.Value.(string)
		}
		if val.Key == "Genre" {
			genre = val.Value.(string)
		}
		if val.Key == "TotalCopies" {
			totalcopies = val.Value.(int32)
		}
		if val.Key == "CurrentCopies" {
			currentcopies = val.Value.(int32)
		}
	}
	return Book{title, author, genre, totalcopies, currentcopies}
}

func GetBooksByAuthor(client *mongo.Client, author string) ([]Book, error) {

	coll := client.Database("Books").Collection("Cluster0")
	var results []bson.D // result of type

	filter := bson.D{{"Author", author}}
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println("some error in obtaing the book")
		return nil, errors.New("Some error")
	}
	var actual []Book
	for cur.Next(context.TODO()) {
		var elem bson.D
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
		datum := getBookFrom(elem)
		actual = append(actual, datum)
	}

	return actual, nil
}
func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func AddBooks(client *mongo.Client, b *Book) (bad error) {
	// this method is used to add new books to the database.
	coll := client.Database("bookHandler").Collection("Cluster0")

	filter := bson.D{{"Author", b.Author}, {"Title", b.Title}}

	cnt, err := coll.CountDocuments(context.TODO(), filter)
	if cnt > 0 {
		return errors.New("Book already present")
	}

	doc, err := toDoc(*b)
	if err != nil {
		log.Fatal(" Problem in marshalling")
		return errors.New("Some problem in inserting")
	}

	_, err = coll.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal("Some Problem in inserting the document")
	}
	return nil
}
