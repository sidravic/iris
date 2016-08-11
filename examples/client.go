package main

import (
	"github.com/supersid/iris/client"
	"fmt"
	"time"
)

func main() {

	_client := client.Start("tcp://127.0.0.1:5555")

	for {
		fmt.Println(_client.SendMessage("echo", "Hello World"))
		time.Sleep(2 * time.Second)
	}
}
