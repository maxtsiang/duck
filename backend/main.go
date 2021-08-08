package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Duck struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Resource string `json:"resource" bson:"resource"`
	Items []string `json:"items,omitempty" bson:"items,omitempty"`
	Photos []string `json:"photos,omitempty" bson:"photos,omitempty"`
}

type User struct {
	ID string
	duck string
	route string
	steps uint32
	pastRoutes []string
	friends []string // duck ids of friends
}

type Route struct {
	ID string
	stops []Stop
}

type Stop struct {
	ID string
	steps uint32 // number of steps needed to reach this stop
}

type Item struct {
	ID string
	url string
	// type (headwear, body, handheld...etc.)
}

type Photo struct {
	ID string
	url string
	rarity uint8 // rarity score for photos
}

// MongoDB global variables
var (
	mongoURI string = "mongodb+srv://mtsiang:030507mc@cluster0.dkc5m.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	database string = "duck-walk"
	ctx context.Context
	ctxCancel context.CancelFunc
	client *mongo.Client
)

func createDuck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var duck Duck
	json.NewDecoder(r.Body).Decode(&duck)

	collection := client.Database(database).Collection("ducks")
	res, err := collection.InsertOne(ctx, duck)
	if err != nil {
		log.Fatal(err)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func connectToMongoDB(ctx context.Context) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
    err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
    return client
}

func main() {
	// Connect to mongodb
	ctx, ctxCancel = context.WithTimeout(context.Background(), 10*time.Second)
    client = connectToMongoDB(ctx)
	defer ctxCancel()
	defer client.Disconnect(ctx)

	// Mux router
	router := mux.NewRouter()
	router.HandleFunc("/api/ducks", createDuck).Methods("POST")

	// HTTP Server
	log.Fatal(http.ListenAndServe(":8080", router))
}