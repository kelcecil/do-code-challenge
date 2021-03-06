package main

import (
	"os"
	"os/signal"
	"syscall"

	"./server"
)

func main() {
	signals := make(chan os.Signal, 1)
	ready := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go ListenForSignal(signals, ready)

	server.StartServer(false, ready)
}

func ListenForSignal(signals chan os.Signal, ready chan bool) {
	<-signals
	ready <- true
}
