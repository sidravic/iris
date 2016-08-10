package main

import (
	"github.com/supersid/iris/broker"
)

func main() {
	broker.Start("tcp://*:5555")
}
