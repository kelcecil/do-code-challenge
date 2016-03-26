package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	ready := make(chan bool, 1)
	go StartServer(ready)
	defer func() { ready <- true }()
	<-ready
	m.Run()
}

func sendMessage(command string) (string, error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(conn)
	conn.Write([]byte(command))
	defer conn.Close()

	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	response = strings.TrimRight(response, "\n")
	return response, nil
}

func expect(t *testing.T, command string, want string, failureMsg string) {
	response, err := sendMessage(command + "\n")
	if err != nil {
		t.Errorf("Error when talking to server; Error: %v", err)
	}
	if response != want {
		t.Errorf(failureMsg+"Expected: %v, Got %v", want, response)
	}
}

func TestServerIntegration(t *testing.T) {
	expect(t, "INDEX|golang|", "OK", "Failed to index first dependency free example. ")
}
