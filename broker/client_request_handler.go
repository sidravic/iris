package broker

import "github.com/supersid/iris/message"

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
	new_service.ProcessRequests()
}