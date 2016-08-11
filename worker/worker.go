package worker

/*
   WORKER_READY
   0: Blank Frame
   1: Command
   2: Blank Frame
   2: Service Name
   3: Unique Worker Identity
*/
import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	uuid "github.com/satori/go.uuid"
//	"github.com/supersid/iris/message"
)

const (
	WORKER_READY    = "WORKER_READY"
	WORKER_RESPONSE = "WORKER_RESPONSE"
)

type Worker struct {
	socket       *zmq.Socket
	identity     string
	service_name string
	broker_url   string
}

type WorkerMessage struct {
	sender          string
	Command         string
	RequestMessage  string
	ResponseMessage string
}

func (worker *Worker) createMessage(command string) []string {
	msg := make([]string, 5)
	msg[0] = ""
	msg[1] = command
	msg[2] = worker.service_name
	msg[3] = worker.identity

	return msg
}

func newWorker(broker_url, service_name string) (*Worker, error) {
	worker := &Worker{
		broker_url:   broker_url,
		service_name: service_name,
	}

	socket, err := zmq.NewSocket(zmq.DEALER)

	if err != nil {
		fmt.Println("[ERROR]: Could not create a new worker socket.")
		return nil, err
	}

	worker.socket = socket
	worker.identity = uuid.NewV4().String()
	poller := zmq.NewPoller()
	poller.Add(worker.socket, zmq.POLLIN)

	return worker, err
}

func (worker *Worker) Close() {
	worker.socket.Close()
}

func (worker *Worker) SendReadyToAccept() ([]string, error) {
	ready_to_accept_msg := worker.createMessage(WORKER_READY)
	_, err := worker.socket.SendMessage(ready_to_accept_msg)

	return ready_to_accept_msg, err
}

func (worker *Worker) wrapMessage(recv_msg []string) WorkerMessage{
	for index, m := range recv_msg {
		fmt.Println(fmt.Sprintf("%d. %s", index, m))
	}

	worker_message := WorkerMessage{sender:        recv_msg[0],
		      			Command:       recv_msg[1],
		      			RequestMessage:recv_msg[2],
	}

	return worker_message
}

func (worker *Worker) unwrapMessage(m WorkerMessage) []string {
	new_message   := make([]string, 5)
	new_message[0] = ""
	new_message[1] = m.sender
	new_message[2] = WORKER_RESPONSE
	new_message[3] = m.RequestMessage
	new_message[4] = m.ResponseMessage
	return new_message
}

func (worker *Worker) Process(m chan WorkerMessage) {
	for {
		msg, err := worker.SendReadyToAccept()

		if err != nil {
			fmt.Println("[ERROR]: Unable to send worker ready message: %s", err.Error())
			fmt.Println(msg)
		}

		received_msg, err := worker.socket.RecvMessage(0)
		worker_message := worker.wrapMessage(received_msg)
		m <- worker_message
	}
}

func (worker *Worker) SendMessage(m WorkerMessage) {
	msg := worker.unwrapMessage(m)
	worker.socket.SendMessage(msg)
}

func Start(broker_url, service_name string) (*Worker, chan WorkerMessage) {
	worker, err := newWorker(broker_url, service_name)

	if err != nil {
		fmt.Println(fmt.Sprintf("[ERROR]: Worker creation error %s", err.Error()))
		panic(err)
	}

	err = worker.socket.Connect(broker_url)

	if err != nil {
		fmt.Println(fmt.Sprintf("[ERROR]: Worker could not connect to broker. %s", err.Error()))
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Starting worker id %s by connecting to %s", worker.identity, worker.broker_url))
	m := make(chan WorkerMessage)
	go worker.Process(m)
	return worker, m
}
