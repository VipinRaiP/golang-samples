// client.go
package main

import (
	"fmt"
	"net"
	"os"
	//"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:12346")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	client_id := os.Args[1]
	message := ""+client_id
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending:", err)
		return
	}

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving:", err)
			return
		}
		response := buffer[:n]
		fmt.Printf("Server response: %s\n", response)
		ack := []byte("ack") 
		_, err = conn.Write(ack)
		if err != nil {
			fmt.Println("Error sending ack:", err)
			return
		}
		fmt.Println("Ack sent")
		//time.Sleep(2)
	}
}
