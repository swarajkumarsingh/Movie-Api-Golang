package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	model "mongoapi/models"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://localhost:27017/netflix"
const dbName = "netflix"
const collectionName = "WatchList"

// MOST IMPORTANT
var collection *mongo.Collection

// Connect to MongoDB
func init() {
	// Client option
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Mongodb")

	collection = client.Database(dbName).Collection(collectionName)

	// Collection instance
	fmt.Println("Collection instance is ready")
}

// Mongodb helpers - file

// Insert 1 record

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted Data", inserted)
	fmt.Println("Inserted Data ID", inserted.InsertedID)
}

func updateOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)

	if err != nil {
		panic(err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		panic(err)
	}

	fmt.Println("Modified count: ", result.ModifiedCount)
}

func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		panic(err)
	}

	fmt.Println("Delete count: ", result.DeletedCount)
}

func deleteAllMovies() int64 {

	filter := bson.D{{}}
	result, err := collection.DeleteMany(context.Background(), filter, nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("Number of movies : ", result.DeletedCount)
	fmt.Println("Number of movies deleted : ", result.DeletedCount)

	return result.DeletedCount
}

func getAllMovies() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	defer cur.Close(context.Background())

	if err != nil {
		panic(err)
	}

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		if err != nil {
			panic(err)
		}

		movies = append(movies, movie)

	}

	return movies

}

func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Alow-Control-Alow-Methods", "POST")
	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)

	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Alow-Control-Alow-Methods", "POST")

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Alow-Control-Alow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Alow-Control-Alow-Methods", "DELETE")

	count := deleteAllMovies()
	json.NewEncoder(w).Encode(count)
}