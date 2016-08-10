package broker

import (
	"fmt"
	"github.com/supersid/iris/service"
	"github.com/supersid/iris/client"
	"github.com/supersid/iris/message"
	"github.com/satori/go.uuid"
)

/*
Data obtained when listening on a ZMQ socket
It's just an slice of strings.
*/

const WORKER_READY string = "WORKER_READY"

//type Message struct {
//	sender       string
//	command      string
//	identity     string
//	data         string
//	service_name string
//}

func (broker *Broker) ParseMessage(msg []string) message.Message {
	var m message.Message
	for index, message := range msg {
		fmt.Println(fmt.Sprintf("%d. %s", index, message))
	}

	command := msg[2]
	sender  := msg[0]

	if command == WORKER_READY || command == client.CLIENT_REQUEST{
		m = message.Message{
			Sender:       sender,
			Command:      command,
			Identity:     msg[4],
			ServiceName:  msg[3],
			Data:         msg[5],
			MessageId:    uuid.NewV4().String(),
		}
	} else {
		m = message.Message{}
	}

	return m
}

func (broker *Broker) ProcessMessage(msg message.Message) {
	fmt.Println("==============================")
	fmt.Println(msg)
	fmt.Println("==============================")
	if msg.Command == WORKER_READY {
		broker.WorkerReadyHandler(msg)
	} else if msg.Command == client.CLIENT_REQUEST {
		fmt.Println("Client Request arrived.")
		broker.ClientRequestHandler(msg)
	}

}

func (broker *Broker) FindOrCreateService(service_name string) (*service.Service, bool) {
	srvc, present := broker.services[service_name]

	if !present {
		srvc = service.NewService(service_name)
	}

	return srvc, present
}
