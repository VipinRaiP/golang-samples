package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"sync"
)

var wg sync.WaitGroup
var count int = 0;

func inc(client *redis.Client) {
	ctx := context.Background()
	done := false
	wg.Add(1)
	for !done{
		val, err := client.SetNX(ctx,"lock", "locked", 0).Result()
		fmt.Println("key was present "+ fmt.Sprintf("%v", val))
		if err != nil {
			fmt.Println("Error !!!")
			panic(err)
			done = true
		} else if val{
			count++;
			client.Del(ctx,"lock")
			done = true		
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	wg.Done()
}

func doInc(client *redis.Client){
	for i:=0;i<10;i++{
		go inc(client)
	}
}

func main() {
	startTime := time.Now()
    client := redis.NewClient(&redis.Options{
        Addr:	  "localhost:6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })

	// ctx := context.Background()

	// err := client.Set(ctx, "foo", "bar", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := client.Get(ctx, "foo").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("foo", val)
	doInc(client)
	wg.Wait()
	fmt.Println(count)	
	timeTaken := time.Since(startTime)
	fmt.Println("Time Taken : ",timeTaken)
}