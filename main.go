package main

func main() {
	ready := make(chan bool)
	StartServer(false, ready)
}
