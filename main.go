package main

import (
	"./network"
	"time"
)


func main() {

	tcpServer := network.NewTCPServer()
	tcpServer.Start(0,0)

	for{
		time.Sleep(1)
	}
}
