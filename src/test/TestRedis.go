package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"client"
)

//func main(){
//	c,err:=redis.Dial("tcp","180.97.69.211:6379")
//	if err!=nil{
//		fmt.Println("connection redis database error:",err)
//	}
//	defer c.Close()

//	c.Do("SET","username","nick")


//	username,err:=redis.String(c.Do("GET","username"))
//	fmt.Println("username:",username)
//}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleClient() {
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exists
}
