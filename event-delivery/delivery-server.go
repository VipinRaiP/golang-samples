// server.go
package main

import (
	"fmt"
	"net"
	"context"
    "github.com/go-redis/redis/v8"
	"time"
	"strconv"
	"encoding/json"
)

type Event struct {
    UserId string `json:"UserId"`
    Payload  string `json:"Payload"`
}


func handleConnection(conn net.Conn, ctx context.Context, client *redis.Client) {
	defer conn.Close()
	// Handle incoming data from the client
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	client_id := string(buffer[:n])
	fmt.Printf("Received: %s\n", client_id)
	fmt.Println("Handling connection for client_id : "+client_id)

	for {
		fmt.Println("Polling redis")
		if(is_there_new_event(client_id,ctx,client)){
			event := get_next_event(client_id,ctx,client)
			// Send a response back to the client
			response,_ := json.Marshal(event)
			_, err = conn.Write(response)
			if err != nil {
				fmt.Println("Error writing:", err)
				return
				// add retry with backoff
			}
			commit_offset(client_id,ctx,client)
		}
		time.Sleep(2*time.Second)
	}
}

func is_there_new_event(client_id string,ctx context.Context,client *redis.Client) bool{
	fmt.Println("Is there new event for : "+client_id)
	offsetStr, err := client.HGet(ctx, "offset",client_id).Result()
	offset,_ := strconv.Atoi(offsetStr)
	if err == redis.Nil {
		// Key does not exist, return the default value
		client.HSet(ctx,"offset",client_id,"-1")
		offset = -1
	} 
	
	len, err := client.LLen(ctx, "events").Result()
	if err!=nil{
		return false
	}
	return (offset+1)<int(len)
}

func get_next_event(client_id string,ctx context.Context,client *redis.Client) Event{
	offset := get_client_offset(client_id,ctx,client)
	offset = offset+1
	offsetStr := strconv.Itoa(offset)
	fmt.Println("Getting next event for client_id "+client_id+" from offset "+offsetStr)
	element, _ := client.LIndex(ctx, "events", int64(offset)).Result()
	fmt.Println("Event retrieved : ", element)
	var event Event
	json.Unmarshal([]byte(element),&event)
	return event
}

func get_client_offset(client_id string,ctx context.Context,client *redis.Client) int{
	fmt.Println("Get client offset : "+client_id)
	offsetStr,_ := client.HGet(ctx, "offset",client_id).Result()
	offset,_ := strconv.Atoi(offsetStr)
	return offset
}

func commit_offset(client_id string, ctx context.Context, client *redis.Client) {
	fmt.Println("Commiting the offset for client : "+client_id)
	offsetStr,_ := client.HGet(ctx,"offset",client_id).Result()
	offset,_ := strconv.Atoi(offsetStr)
	offset = offset + 1
	offsetStr = strconv.Itoa(offset)
	client.HSet(ctx,"offset",client_id,offsetStr)
}

func main() {

	ctx := context.Background()

    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Replace with your Redis server's address
        Password: "",          // No password by default
        DB: 0,                 // Default DB
    })

    defer client.Close()

	listener, err := net.Listen("tcp", "localhost:12346")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on localhost:12346")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go handleConnection(conn,ctx,client)
	}
}
