package broker

import (
	"fmt"
	"github.com/supersid/iris/service"
)

/*
Data obtained when listening on a ZMQ socket
It's just an slice of strings.
*/

const WORKER_READY string = "WORKER_READY"

type Message struct {
	sender       string
	command      string
	identity     string
	data         string
	service_name string
}

func (broker *Broker) ParseMessage(msg []string) Message {
	var message Message
	for index, message := range msg {
		fmt.Println(fmt.Sprintf("%d. %s", index, message))
	}

	command := msg[2]
	sender := msg[0]

	if command == WORKER_READY {
		message = Message{
			sender:       sender,
			command:      command,
			identity:     msg[4],
			service_name: msg[3],
			data:         "",
		}
	} else {
		message = Message{}
	}

	return message
}

func (broker *Broker) ProcessMessage(msg Message) {
	if msg.command == WORKER_READY {
		new_service := broker.FindOrCreateService(msg.service_name)
		worker_existed, service_worker := new_service.FindOrCreateServiceWorker(msg.identity, msg.sender)

		if !worker_existed {
			fmt.Println("Adding service worker")
			new_service.AddWorker(service_worker)
		}

		fmt.Println("Service after workers are added")

		fmt.Println(len(new_service.GetWaitingWorkers()))
		for index, w := range new_service.GetWaitingWorkers() {
			fmt.Println(fmt.Sprintf("%d. %s", index, w.GetIdentity()))
		}
	}
}

func (broker *Broker) FindOrCreateService(service_name string) (*service.Service) {
	srvc, present := broker.services[service_name]

	if !present {
		srvc = service.NewService(service_name)
		broker.services[service_name] = srvc
		broker.services_list = append(broker.services_list, srvc)
	}

	return srvc
}
