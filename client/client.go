package client

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	uuid "github.com/satori/go.uuid"
	"errors"
)

const (
	CLIENT_REQUEST = "CLIENT_REQUEST"
)

type Client struct {
	socket *zmq.Socket
	broker_url string
	identity   string
}

func (client *Client) setIdentity() {
	client.identity = uuid.NewV4().String()
}

func newClient(broker_url string) (*Client, error) {
	client := &Client{
		broker_url: broker_url,
	}

	socket, err := zmq.NewSocket(zmq.DEALER)

	if err != nil {
		fmt.Println("[ERROR]: Could not create a new client socket.")
		return nil, err
	}

	client.socket = socket
	client.setIdentity()
	return client, err
}

func (client *Client) Close() {
	client.socket.Close()
}

/*
   CLIENT_REQUEST
   0: Blank Frame
   1: Command
   2: Blank Frame
   2: Service Name
   3: Data
*/
func (client *Client) createMessage(service_name, message string) []string {
	msg := make([]string, 5)
	msg[0] = ""
	msg[1] = CLIENT_REQUEST
	msg[2] = ""
	msg[3] = service_name
	msg[4] = message

	return msg
 }

func (client *Client) SendMessage(service_name, message string) error{
	if service_name == "" {
		return errors.New("service_name cannot be nil or blank.")
	}

	msg := client.createMessage(service_name, message)
	fmt.Println(fmt.Sprintf("Sending %s: to service %s", message, service_name))
	fmt.Println(fmt.Sprintf("Encoded message :%s", msg))
	_, err := client.socket.SendMessage(msg)

	return err
}

func Start(broker_url string) *Client{
	client, err := newClient(broker_url)

	if err != nil {
		fmt.Println(fmt.Sprintf("[ERROR]: Client creation error %s", err.Error()))
		panic(err)
	}

	err = client.socket.Connect(broker_url)

	if err != nil {
		fmt.Println(fmt.Sprintf("[ERROR]: Client could not connect to broker. %s", err.Error()))
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Starting client id %s by connecting to %s", client.identity, client.broker_url))
	return client
}
