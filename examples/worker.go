package main

import ()
import "github.com/supersid/iris/worker"

func main() {
	worker.Start("tcp://127.0.0.1:5555", "echo")
}
