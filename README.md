# IRIS

Iris is a service discovery and load balancing library built on top of 
ZMQ. It uses the Majordomo Pattern to build a basic DEALER-ROUTER-DEALER
 pattern. 
 
It does not adhere to the exact protocol specified in the zmq specifications 
for message exchange. 

Basic examples are available in the `examples` directory.

## Setup 

After pulling dependencies using Godep and installing Zeromq 4.1 launch each of the components
on separate terminals

1. `go run examples/broker.go`
2. `go run examples/client.go`
3. `go run examples/worker.go`
4. `go run examples/worker.go`
5. `go run examples/worker.go`

You may run as many copies of the worker and the client. 

## TODO

1. Heartbeat checks
2. Service Listings
3. Realtime monitoring
4. Logging of events.


