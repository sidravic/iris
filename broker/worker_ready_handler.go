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
	new_service, already_present := broker.FindOrCreateService(msg.ServiceName)

	if !already_present {
		broker.AddService(new_service)
	}

	worker_existed, service_worker := new_service.FindOrCreateServiceWorker(msg.Identity, msg.Sender)

	if !worker_existed {
		fmt.Println("Adding service worker")
		new_service.AddWorker(service_worker)
		broker.addWorker(new_service, service_worker)
	}

	for index, w := range new_service.GetWaitingWorkers() {
		fmt.Println(fmt.Sprintf("%d. %s", index, w.GetIdentity()))
	}
}

func (broker *Broker) AddService(service *service.Service) {
	service_name := service.GetName()
	broker.services[service_name] = service
	broker.services_list = append(broker.services_list, service)
}

func (broker *Broker) addWorker(service *service.Service,
				service_worker *service.ServiceWorker){
	service_name := service.GetName()
	broker.workers[service_name] = service_worker
}