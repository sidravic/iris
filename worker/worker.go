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
)

const (
	WORKER_READY = "WORKER_READY"
)

type Worker struct {
	socket       *zmq.Socket
	identity     string
	service_name string
	broker_url   string
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

func (worker *Worker) Process() {
	for {
		msg, err := worker.SendReadyToAccept()

		if err != nil {
			fmt.Println("[ERROR]: Unable to send worker ready message: %s", err.Error())
			fmt.Println(msg)
		}

		received_msg, err := worker.socket.RecvMessage(0)

		for index, message := range received_msg {
			fmt.Println(fmt.Sprintf("%d. %s", index, message))
		}
	}
}

func Start(broker_url, service_name string) {
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

	defer worker.Close()
	fmt.Println(fmt.Sprintf("Starting worker id %s by connecting to %s", worker.identity, worker.broker_url))
	worker.Process()
}
