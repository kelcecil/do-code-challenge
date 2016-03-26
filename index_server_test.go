package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

func TestServerIntegration(t *testing.T) {
	ready := make(chan bool, 1)
	go StartServer(ready)
	//defer func() { quit <- true }()

	<-ready
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Errorf("Failed to connect to local server; Error: %v", err)
		return
	}

	reader := bufio.NewReader(conn)

	indexGolang := "INDEX|golang|\n"
	conn.Write([]byte(indexGolang))
	response, err := reader.ReadString('\n')
	response = strings.TrimRight(response, "\n")

	if err != nil {
		t.Errorf("Failed to read response from local server; Error: %v", err)
	}
	if response != "OK" {
		t.Errorf("Failed to index first dependency free example")
	}

}
