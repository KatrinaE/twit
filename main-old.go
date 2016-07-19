package main

import(
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
	"os"
)

func main() {
	client, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("uh oh, error connecting to redis")
		os.Exit(1)
	}
	fmt.Println("Connected!")

	foo, err := client.Cmd("GET", "foo").Str()
	if err != nil {
		fmt.Println("error getting foo")
	}
	fmt.Printf("Foo is: %s\n", foo)
	
}
