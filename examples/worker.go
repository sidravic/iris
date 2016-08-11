package main

import (
	"github.com/supersid/iris/worker"
	"fmt"
	"time"
)

func main() {
	w, c := worker.Start("tcp://127.0.0.1:5555", "echo")

	fmt.Println(w)
	defer w.Close()

	for {
		fmt.Println("From Message")

		msg := <- c
		fmt.Println(fmt.Sprintf("%s", msg.RequestMessage))
		msg.ResponseMessage = fmt.Sprintf("The answer is %s", time.Now())
		w.SendMessage(msg)
	}


}
