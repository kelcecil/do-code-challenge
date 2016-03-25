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

	brokenSampleOne  string   = "INDEX|cloog|bloob|gmp,isl,pkg-config\n"
	brokenMessageOne *Message = NewMessage("INDEX", "cloog", []string{"gmp", "isl", "pkg-config"})

	brokenSampleTwo  string   = "INDEX|cloog,apples|gmp,isl,pkg-config\n"
	brokenMessageTwo *Message = NewMessage("INDEX", "cloog", []string{"gmp", "isl", "pkg-config"})
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

func BenchmarkParseMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseMessage(sampleOne)
		ParseMessage(sampleTwo)
		ParseMessage(sampleThree)
		ParseMessage(sampleFour)
	}
}

func TestSampleBrokenMessages(t *testing.T) {
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

	testMessage(brokenSampleOne, brokenMessageOne, "INDEX message with extra bars failed.")
}
