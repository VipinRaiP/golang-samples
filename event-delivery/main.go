package main

import (
    "context"
    "fmt"
    "os"
    "strconv"
	"encoding/json"
    "github.com/go-redis/redis/v8"
	"strconv"
)

type Event struct {
    UserId string `json:"UserId"`
    Payload  string `json:"Payload"`
}


func main() {
    ctx := context.Background()

    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Replace with your Redis server's address
        Password: "",          // No password by default
        DB: 0,                 // Default DB
    })

    defer client.Close()

    // parse command

	go listen_for_user(ctx,client)

	// poll redis and send to user

	go poll_and_deliver(ctx,client)


	// if len(os.Args) != 3 {
	// 	fmt.Println("Usage: go run main.go <key> <index>")
	// 	return
	// }

	// key := os.Args[1]
	// index, err := strconv.Atoi(os.Args[2])
	// if err != nil {
	// 	fmt.Println("Invalid index:", err)
	// 	return
	// }
	// event := get_from_list(key, index, ctx, client)
	// fmt.Printf("%+v\n",event)
}

func listen_for_user(ctx context.Context, client *redis.Client){
	for {
		fmt.Println("listening for clients..")
		command := os.Args[1]
		if command=="listen"{
			user_id, _ := strconv.Atoi(os.Args[2])
			go handle_user(user_id,client,ctx)
		}else{
			fmt.Println("Incorrect command retry...")
		}
	}	 
}

func handle_user(user_id int, ctx context.Context, client *redis.Client){
	key := strconv.Itoa(user_id)
	field := "is_online"
	value := true
	defer mark_offline(user_id, ctx, client)
	fmt.Println("Marking user with id "+key+"as online")
	err := client.HSet(ctx, key, field, value).Err()
    if err != nil {
        fmt.Println("Setting status failed :", err)
    }
	
	
}

func mark_offline(user_id int, ctx context.Context, client *redis.Client){
	key := strconv.Itoa(user_id)
	field := "is_online"
	value := false
	fmt.Println("Marking user with id "+key+"as offline")
	err := client.HSet(ctx, key, field, value).Err()
    if err != nil {
        fmt.Println("Setting status failed :", err) 
    }
}


func get_from_list(key string, index int, ctx context.Context, client *redis.Client) Event {
	value, err := client.LIndex(ctx, key, int64(index)).Result()
    var event Event
	if err != nil {
        fmt.Println("Failed to read from Redis:", err)
        return event
    }
	err = json.Unmarshal([]byte(value),&event)
	if err != nil{
		fmt.Println("Json unmarshaling error: ",err)
		return event
	}
	fmt.Printf("Value string at index %d: %s\n", index, value)
    fmt.Printf("Value at index %d: %+v\n", index, event)
	return event
}


// TODOS


// to do regist client
// get client hearbeat
// poll and deliver
// handle failures
