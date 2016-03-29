package server

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"../message"
	"../parser"
)

var (
	activeConnections sync.WaitGroup = sync.WaitGroup{}
)

// Test mode and ready is to let the integration tests
// know that we're ready to begin.
func StartServer(testMode bool, ready chan bool) {

	msgRouter := make(chan *message.Message, 100)
	go message.MessageRouter(msgRouter)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	log.Print("Server has started.")

	// If this part of an integration test, then
	// let the test know that it can continue.
	if testMode {
		ready <- true
	}

	for {
		select {
		case <-ready:
			log.Print("Received SIGTERM. Qutting...")
			listener.Close()
			activeConnections.Wait()
			return
		default:
		}

		// Set a deadline for accept so that we can check for signals
		if listener, ok := listener.(*net.TCPListener); ok {
			listener.SetDeadline(time.Now().Add(30 * time.Second))
		}

		connection, err := listener.Accept()
		if err != nil {
			if err.(*net.OpError).Timeout() {
				continue
			}
		}
		activeConnections.Add(1)

		go HandleConnection(connection, msgRouter)
	}
}

func HandleConnection(conn net.Conn, msgRouter chan<- *message.Message) {
	log.Printf("Connection open.")
	reader := parser.NewMessageReader(conn)
	for {
		conn.SetDeadline(time.Now().Add(30 * time.Second))
		message, err := reader.Read()
		if err == io.EOF {
			log.Printf("Closing connection")
			conn.Close()
			activeConnections.Done()
			break
		} else if err != nil {
			_, err = conn.Write([]byte("ERROR\n"))
			if err != nil {
				log.Print(err)
			}
			continue
		}

		msgRouter <- message
		response := <-message.Response
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			log.Print(err)
		}
	}
}
