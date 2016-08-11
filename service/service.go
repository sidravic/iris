package service

import (
	"time"
	"github.com/supersid/iris/message"
	"fmt"
	"errors"
)

type Service struct {
	name string
	waitingWorkers []*ServiceWorker
	lastHeartbeat  time.Time
	waitingRequests []message.Message
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
		waitingWorkers:  make([]*ServiceWorker, 0),
		waitingRequests: make([]message.Message, 0),
	}
	return service
}

func (service *Service) GetName() string{
	return service.name
}

func (service *Service) GetWaitingWorkers() []*ServiceWorker{
	return service.waitingWorkers
}

func (service *Service) FindOrCreateServiceWorker(identity string, sender string) (bool, *ServiceWorker) {
	var service_worker *ServiceWorker
	var worker_exists bool = false

	for i := 0; i < len(service.waitingWorkers); i++ {
		if service.waitingWorkers[i].identity == identity {
			worker_exists = true
			service_worker = service.waitingWorkers[i]
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
	service.waitingWorkers = append(service.waitingWorkers, service_worker)
}

/*
Add message to services.waiting_requests if a message (read request) with the
same message Id is not already present.
 */
func (service *Service) AddRequest(msg message.Message){
	message_exists := false

	for _, request := range service.waitingRequests {
		if request.MessageId == msg.MessageId {
			message_exists = true
			break;
		}
	}

	if !message_exists {
		service.waitingRequests = append(service.waitingRequests, msg)
	}
}

func (service *Service) ProcessRequests() (error, message.Message, *ServiceWorker) {
	var err error
	if len(service.waitingWorkers) <= 0 {
		fmt.Println("No workers available at the moment")
		err = errors.New("No workers available at the moment")

	}

	if len(service.waitingRequests) <= 0 {
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
	request := service.waitingRequests[0]
	service.waitingRequests = service.waitingRequests[1:]
	return request
}

func (service *Service) popFirstWorker() *ServiceWorker{
	service_worker := service.waitingWorkers[0]
	service.waitingWorkers = service.waitingWorkers[1:]
	return service_worker
}

