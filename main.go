package main

func main() {
	quit := make(chan bool)
	StartServer(quit)
}
