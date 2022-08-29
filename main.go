package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jainachal03/lib/handlers"
	"github.com/jainachal03/lib/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createConnection(username, password string) (client *mongo.Client, err error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://" + username + ":" + password + "@cluster0.byopwfw.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client, nil
}
func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	client, err := createConnection(config.DB_USERNAME, config.DB_PASSWORD)

	if err != nil {
		log.Fatal(err)
	}
	l := log.New(os.Stdout, "dat", log.LstdFlags)
	bh := handlers.NewBooks(l, client)
	ah := handlers.NewAuthHandler(l, client)
	s := mux.NewRouter()

	getRouter := s.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", bh.GetBooks)

	getRouter.HandleFunc("/book/{Author}", bh.GetBooksByAuthor)

	postRouter := s.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/add", bh.AddBooks)
	postRouter.HandleFunc("/auth", ah.Login)
	postRouter.HandleFunc("/signup", ah.Signup)

	http.ListenAndServe(":8080", s)
	// coll := client.Database("bookHandler").Collection("Cluster0")

	// filter := bson.D{{"Author", "Jane Austen"}}
	// cursor, err := coll.Find(context.TODO(), filter)
	// if err != nil {
	// 	panic(err)
	// }
	// var results []bson.D
	// if err = cursor.All(context.TODO(), &results); err != nil {
	// 	panic(err)
	// }
	// for _, result := range results {
	// 	fmt.Println(result)
	// }

}
