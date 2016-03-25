package main

import (
	"fmt"
	"testing"
)

var (
	sampleOne  string   = "INDEX|cloog|gmp,isl,pkg-config\n"
	messageOne *Message = NewMessage("INDEX", "cloog", []string{"gmp", "isl", "pkg-config"})

	sampleTwo  string   = "INDEX|ceylon|\n"
	messageTwo *Message = NewMessage("INDEX", "ceylon", []string{})

	sampleThree  string   = "REMOVE|cloog|\n"
	messageThree *Message = NewMessage("REMOVE", "cloog", []string{})

	sampleFour  string   = "QUERY|cloog|\n"
	messageFour *Message = NewMessage("QUERY", "cloog", []string{})
)

func TestSampleFullMessages(t *testing.T) {
	testMessage := func(sample string, message *Message, failureMessage string) {
		parsedMessage, err := ParseMessage(sample)
		if err != nil {
			t.Errorf("Parsing of message %v failed because: %v", sample, err.Error())
		}
		if !message.Equals(parsedMessage) {
			message := fmt.Sprintf("%v -- %v", failureMessage, parsedMessage.String())
			t.Errorf(message)
		}
	}

	testMessage(sampleOne, messageOne, "INDEX message with dependencies failed.")
	testMessage(sampleTwo, messageTwo, "INDEX message without dependencies failed.")
	testMessage(sampleThree, messageThree, "REMOVE message failed.")
	testMessage(sampleFour, messageFour, "QUERY message failed.")
}
