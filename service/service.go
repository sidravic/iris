package service

import (
	"time"
	"github.com/supersid/iris/message"
	"fmt"
	"errors"
)

type Service struct {
	name string
	waiting_workers []*ServiceWorker
	last_heartbeat  time.Time
	waiting_requests []message.Message
}

type ServiceWorker struct {
	sender       string
	identity     string
}

func (service_worker *ServiceWorker) GetIdentity() string {
	return service_worker.identity
}

func (service_worker *ServiceWorker) GetSender() string {
	return service_worker.sender
}

func NewService(name string) (*Service){
	service := &Service{
		name:             name,
		waiting_workers:  make([]*ServiceWorker, 0),
		waiting_requests: make([]message.Message, 0),
	}
	return service
}

func (service *Service) GetName() string{
	return service.name
}

func (service *Service) GetWaitingWorkers() []*ServiceWorker{
	return service.waiting_workers
}

func (service *Service) FindOrCreateServiceWorker(identity string, sender string) (bool, *ServiceWorker) {
	var service_worker *ServiceWorker
	var worker_exists bool = false

	for i := 0; i < len(service.waiting_workers); i++ {
		if service.waiting_workers[i].identity == identity {
			worker_exists = true
			service_worker = service.waiting_workers[i]
			break;
		}
	}

	if !worker_exists {
		service_worker = &ServiceWorker{
			sender:   sender,
			identity: identity,
		}
	}

	return worker_exists, service_worker
}


func (service *Service) AddWorker(service_worker *ServiceWorker) {
	service.waiting_workers = append(service.waiting_workers, service_worker)
}

/*
Add message to services.waiting_requests if a message (read request) with the
same message Id is not already present.
 */
func (service *Service) AddRequest(msg message.Message){
	message_exists := false

	for _, request := range service.waiting_requests {
		if request.MessageId == msg.MessageId {
			message_exists = true
			break;
		}
	}

	if !message_exists {
		service.waiting_requests = append(service.waiting_requests, msg)
	}
}

func (service *Service) ProcessRequests() (error, message.Message, *ServiceWorker) {
	var err error
	if len(service.waiting_workers) <= 0 {
		fmt.Println("No workers available at the moment")
		err = errors.New("No workers available at the moment")

	}

	if len(service.waiting_requests) <= 0 {
		fmt.Println("No requests available to process at the moment")
		err = errors.New("No requests available to process at the moment")
	}

	if err != nil {
		return err, message.Message{}, &ServiceWorker{}
	}

	msg := service.popFirstRequest()
	service_worker := service.popFirstWorker()
	return err, msg, service_worker

}


func (service *Service) popFirstRequest() message.Message {
	request := service.waiting_requests[0]
	service.waiting_requests = service.waiting_requests[1:]
	return request
}

func (service *Service) popFirstWorker() *ServiceWorker{
	service_worker := service.waiting_workers[0]
	service.waiting_workers = service.waiting_workers[1:]
	return service_worker
}

