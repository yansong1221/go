package main

import (
	"dispatch"
	"network"
	"time"
)

func main() {

	eventDispatch := dispatch.NewEventDispatch()
	tcpServer := network.NewTCPServer(eventDispatch, 1024)

	tcpServer.Start(8888)

	for {
		eventDispatch.Update()
		time.Sleep(1)
	}
}
