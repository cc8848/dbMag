package main

import (
	"fmt"
	"github.com/go-redis/redis"
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

func GetRedisConn(redisHost string, passwd string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: passwd, // no password set
		DB:       0,      // use default DB
	})
	return client
}
func main() {
	var conn *redis.Client
	addr := "localhost:6379"
	password := ""
	conn = GetRedisConn(addr, password)

	pong, err := conn.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	key := "username"
	value := "ghc"
	conn.Set(key, value, 0)
}

//http://www.jianshu.com/p/7e22ad3a9061
//http://siddontang.com/categories/go/
http://faemalia.com/Technology/mysql-internals.pdf
func ExampleClient(client *redis.Client) {
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
