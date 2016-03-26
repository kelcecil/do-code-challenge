package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

var (
	packages *PackageSet = NewPackageSet()
)

func StartServer(ready chan bool) {
	log.Print("Starting server")

	msgRouter := make(chan *Message, 500)
	go MessageRouter(msgRouter)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	ready <- true

	for {
		select {
		case <-ready:
			listener.Close()
			return
		default:
		}

		log.Print("Waiting for connection")
		connection, err := listener.Accept()

		log.Print("Accepted connection")
		if err != nil {
			panic(err)
		}

		go HandleConnection(connection, msgRouter)

	}
}

func HandleConnection(conn net.Conn, msgRouter chan<- *Message) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Print("Connection closed")
			conn.Close()
			return
		} else if err != nil {
			log.Print(err)
		}
		message, _ := ParseMessage(line)
		msgRouter <- message
		log.Print("Reading response")
		response := <-message.Response
		log.Print(response)
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			log.Print(err)
		}
	}
}
