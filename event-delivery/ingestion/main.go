package main

import (
    "fmt"
    "net/http"
	"encoding/json"
    //"github.com/gorilla/mux"
	"context"
	"github.com/go-redis/redis/v8"
)

type Event struct {
    UserId string `json:"UserId"`
    Payload  string `json:"Payload"`
}

func HelloHandler(client *redis.Client, ctx context.Context) http.HandlerFunc{
    return func(w http.ResponseWriter, r *http.Request){
		// Parse JSON from the request body
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "Failed to decode JSON payload", http.StatusBadRequest)
			return
		}
	
		// Do something with the event data
		fmt.Println("Received event:", event)


		key := "events" // Replace with your desired key name
	
		jsonBytes, err := json.Marshal(event)
		if err != nil {
			fmt.Println("JSON marshaling error:", err)
			return
		}

		result, err := client.LPush(ctx, key, jsonBytes).Result()
		if err != nil {
			fmt.Println("Failed to push JSON to Redis:", err)
			return
		}
	
		fmt.Printf("Added to list at index %d\n", result-1)
		w.WriteHeader(http.StatusOK)
    	fmt.Fprintln(w, "JSON payload received successfully")
	}
}

func main() {
	

	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Replace with your Redis server's address
		Password: "",          // No password by default
		DB: 0,                 // Default DB
	})

	http.HandleFunc("/ingest", HelloHandler(client,ctx))

    http.ListenAndServe(":8080", nil)
	client.close()
}
