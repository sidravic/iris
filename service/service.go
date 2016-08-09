package service

import (
	"time"
)

type Service struct {
	name string
	waiting_workers []*ServiceWorker
	last_heartbeat  time.Time
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
		name: name,
		waiting_workers: make([]*ServiceWorker, 0),
	}
	return service
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