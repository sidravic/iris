package broker

import (
	"github.com/supersid/iris/message"
	"github.com/supersid/iris/service"
)

const WORKER_REQUEST = "WORKER_REQUEST"

func(broker *Broker) ClientRequestHandler(msg message.Message){
	/*
	Find if a service exists
	1. Find or create service
	2. Add service to broker
	3. Assign requests to waiting requests queue
	4.
	*/
	new_service, already_present := broker.FindOrCreateService(msg.ServiceName)

	if !already_present {
		broker.AddService(new_service)
	}

	new_service.AddRequest(msg)
	err, msg, service_worker := new_service.ProcessRequests()

	if err != nil {
		return
	}

	broker.ProcessClientRequest(service_worker, msg)
}


func (broker *Broker) ProcessClientRequest(service_worker *service.ServiceWorker,
				           msg message.Message){
	client_sender := msg.Sender
	new_message := make([]string, 5)
	new_message[0] = service_worker.GetSender()
	new_message[1] = client_sender
	new_message[2] = WORKER_REQUEST
	new_message[3] = msg.Data

	broker.socket.SendMessage(new_message)
}

