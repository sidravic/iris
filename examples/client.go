package main

import (
	"github.com/supersid/iris/client"
	"fmt"
)

func main() {

	_client := client.Start("tcp://127.0.0.1:5555")
	fmt.Println(_client.SendMessage("echo", "Hello World"))
	for {}
}
