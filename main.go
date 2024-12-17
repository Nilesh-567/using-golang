package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
)

var collection *mongo.Collection

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func connectDB() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

func init() {
	client := connectDB()
	collection = client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func signUpHandlers() {
	fmt.Printf("get method responds successfully ")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signup", signUpHandler).Methods("POST")
        //s := mux.NewRouter()
	//s.HandleFunc("/", signUpHandlers).Methods("GET")
	
	port := os.Getenv("PORT")
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
	//log.Fatal(http.ListenAndServe(":"+port, s))
}
