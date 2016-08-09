package client

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	uuid "github.com/satori/go.uuid"
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

func Start(broker_url string) {
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

	defer client.Close()
	fmt.Println(fmt.Sprintf("Starting client id %s by connecting to %s", client.identity, client.broker_url))
}
