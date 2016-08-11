package main

import ()
import (
	"github.com/supersid/iris/worker"
	"fmt"
)

func main() {
	w, c := worker.Start("tcp://127.0.0.1:5555", "echo")

	fmt.Println(w)
	defer w.Close()

	for {
		fmt.Println("From Message")

		msg := <- c
		fmt.Println(fmt.Sprintf("%s", msg.RequestMessage))
		msg.ResponseMessage = "The answer is 42."
		w.SendMessage(msg)
	}


}
