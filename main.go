package main

import "github.com/kelcecil/do-code-challenge/server"

func main() {
	ready := make(chan bool)
	server.StartServer(false, ready)
}
