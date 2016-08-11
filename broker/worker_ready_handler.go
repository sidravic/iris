package broker

import (
	"fmt"
	"github.com/supersid/iris/service"
	"github.com/supersid/iris/message"
)

/*
 When the worker is read
 1. Find or create a service
 2. Add the service to the list of the services the brokers serves
 3. Find or create a new service worker under the service.
 4. Add the worker to the service
 5. Add the worker to the broker.
 */
func (broker *Broker) WorkerReadyHandler(msg message.Message) {
	newService, alreadyPresent := broker.FindOrCreateService(msg.ServiceName)

	if !alreadyPresent {
		broker.AddService(newService)
	}

	workerExisted, serviceWorker := newService.FindOrCreateServiceWorker(msg.Identity, msg.Sender)

	if !workerExisted {
		fmt.Println("Adding service worker")
		newService.AddWorker(serviceWorker)
		broker.addWorker(newService, serviceWorker)
	}

	for index, w := range newService.GetWaitingWorkers() {
		fmt.Println(fmt.Sprintf("%d. %s", index, w.GetIdentity()))
	}
}

func (broker *Broker) AddService(service *service.Service) {
	service_name := service.GetName()
	broker.services[service_name] = service
	broker.servicesList = append(broker.servicesList, service)
}

func (broker *Broker) addWorker(service *service.Service,
				service_worker *service.ServiceWorker){
	service_name := service.GetName()
	broker.workers[service_name] = service_worker
}