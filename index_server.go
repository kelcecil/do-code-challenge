package main

import (
	"bufio"
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

		log.Print("Waiting for connection")
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
		log.Print(line)
		message, err := ParseMessage(line)
		if err != nil {
			_, err = conn.Write([]byte("ERROR\n"))
			continue
		}
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
