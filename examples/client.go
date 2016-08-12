package main

import (
	"github.com/supersid/iris/client"
	"fmt"

	"time"
)

func SendMessages(_client *client.Client){
	seq := 0

	for {
		seq++
		msg := fmt.Sprintf("Hello %d", seq)
		fmt.Println(_client.SendMessage("echo", msg))
		time.Sleep(250 * time.Millisecond)
		fmt.Println(seq)
	}
}

func main() {

	_client := client.Start("tcp://127.0.0.1:5555")
	go SendMessages(_client)

	for {
		err, msg := _client.ReceiveMessage()
		fmt.Println("------------------------")
		fmt.Println(err)
		fmt.Println(msg)
		fmt.Println("------------------------")
	}



}
