package broker

import (
	"github.com/supersid/iris/service"
	"github.com/supersid/iris/worker"
	zmq "github.com/pebbe/zmq4"
	"time"
	"fmt"
)

const POLL_FREQUENCY = 250 * time.Millisecond


type Broker struct {
	socket          *zmq.Socket
	broker_url      string
	services        map[string]*service.Service
	workers         map[string]*worker.Worker
	services_list   []*service.Service
	waiting_workers []*worker.Worker
}

/*
 Creates a new broker by initialising a new ROUTER socket
 */
func newBroker(broker_url string) (*Broker, error) {
	broker := &Broker{
		broker_url:      broker_url,
		services:        make(map[string]*service.Service),
		workers:         make(map[string]*worker.Worker),
		waiting_workers: make([]*worker.Worker,0),
	}

	socket, err := zmq.NewSocket(zmq.ROUTER)

	if err != nil {
		fmt.Println("[ERROR]: Unable to create socket.")
		return nil, err
	}

	broker.socket = socket
	return broker, nil
}

/*
 Launches the poller to periodically check for messages
 */
func (broker *Broker) Process() {
	poller := zmq.NewPoller()
	poller.Add(broker.socket, zmq.POLLIN)

	for {
		fmt.Println("Polling...")
		incoming_sockets, err := poller.Poll(POLL_FREQUENCY)
		if err != nil {
			fmt.Println(err)
		}

		if len(incoming_sockets) > 0 {
			msg, _ := broker.socket.RecvMessage(0)
			var message Message = broker.ParseMessage(msg)
			fmt.Println("Printing parsed message %s", message)
			broker.ProcessMessage(message)
		}

	}
}

func (broker *Broker) GetAllServices() []*service.Service {
	return broker.services_list
}

func (broker *Broker) Close() error {
	err := broker.socket.Close()
	return err
}

func Start(broker_url string) {
	broker, err := newBroker(broker_url)

	if err != nil {
		fmt.Println("[ERROR]: Broker creation failed.")
		panic(err)
	}

	err = broker.socket.Bind(broker_url)

	if err != nil {
		fmt.Println("[ERROR]: Broker Bind failed.")
		panic(err)

	}

	defer broker.Close()
	fmt.Println(fmt.Sprintf("Starting broker on %s", broker.broker_url))
	broker.Process()
}

