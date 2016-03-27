package main

import (
	"io"
	"log"
	"net"
	"time"
)

var (
	packages *PackageSet = NewPackageSet()
)

func StartServer(testMode bool, ready chan bool) {
	log.Print("Starting server")

	msgRouter := make(chan *Message, 1000)
	go MessageRouter(msgRouter)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	if testMode {
		ready <- true
	}

	for {
		select {
		case <-ready:
			listener.Close()
			return
		default:
		}

		if listener, ok := listener.(*net.TCPListener); ok {
			listener.SetDeadline(time.Now().Add(5 * time.Second))
		}

		connection, err := listener.Accept()
		if err != nil {
			if err.(*net.OpError).Timeout() {
				continue
			}
		}
		log.Print("Accepted connection")

		go HandleConnection(connection, msgRouter)
	}
}

func HandleConnection(conn net.Conn, msgRouter chan<- *Message) {
	reader := NewMessageReader(conn)
	for {
		message, err := reader.Read()
		if err == io.EOF {
			log.Print("Connection closed")
			conn.Close()
			break
		} else if err != nil {
			_, err = conn.Write([]byte("ERROR\n"))
			if err != nil {
				log.Print(err)
			}
			continue
		}

		log.Print("Sending msg")
		msgRouter <- message
		log.Print("Waiting for response")
		response := <-message.Response
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			log.Print(err)
		}
	}
}
