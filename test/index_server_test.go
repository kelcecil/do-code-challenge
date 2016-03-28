package test

import (
	"bufio"
	"net"
	"strings"
	"testing"

	"../server"
)

func TestMain(m *testing.M) {
	ready := make(chan bool, 1)
	go server.StartServer(true, ready)
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
		t.Errorf(failureMsg+" Expected: %v, Got %v", want, response)
	}
}

func TestIntegrationScenarioOne(t *testing.T) {
	expect(t, "QUERY|golang|", "FAIL", "An query was OK that should not exist.")
	expect(t, "INDEX|golang|", "OK", "Failed to index first dependency free example.")
	expect(t, "QUERY|golang|", "OK", "A query was FAIL that should exist.")
	expect(t, "INDEX|glide|golang", "OK", "Failed to index a package with a known dependency.")
	expect(t, "QUERY|glide|", "OK", "A package with a dependency was not indexed and should be.")
	expect(t, "REMOVE|golang|", "FAIL", "A package which has a dependent was deleted and should not have been.")
	expect(t, "QUERY|golang|", "OK", "A query was FAIL that should exist.")
	expect(t, "REMOVE|glide|", "OK", "A package with no dependents failed to be removed.")
	expect(t, "REMOVE|golang|", "OK", "A package with no CURRENT dependents failed to be removed.")
}

func TestIntegrationScenarioTwo(t *testing.T) {
	// Broken message
	expect(t, "INDEX|emacs☃elisp", "ERROR", "A query that should have errored did not.")
	expect(t, "WAT|morewat|", "ERROR", "A made up command should return error.")
}

func TestIntegrationScenarioThree(t *testing.T) {
	expect(t, "INDEX|golo|java", "FAIL", "A required dependency doesn't exist.")
	expect(t, "INDEX|java|", "OK", "Failed to insert a package without deps.")
	expect(t, "INDEX|golo|java", "OK", "Failed to insert a pkg after dep is indexed.")
}
