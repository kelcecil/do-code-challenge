package test

import (
	"fmt"
	"testing"

	"../message"
	"../parser"
)

var (
	sampleOne  string           = "INDEX|cloog|gmp,isl,pkg-config\n"
	messageOne *message.Message = message.NewMessage("INDEX", "cloog", []string{"gmp", "isl", "pkg-config"})

	sampleTwo  string           = "INDEX|ceylon|\n"
	messageTwo *message.Message = message.NewMessage("INDEX", "ceylon", []string{})

	sampleThree  string           = "REMOVE|cloog|\n"
	messageThree *message.Message = message.NewMessage("REMOVE", "cloog", []string{})

	sampleFour  string           = "QUERY|cloog|\n"
	messageFour *message.Message = message.NewMessage("QUERY", "cloog", []string{})

	brokenSampleOne   string = "INDEX|emacs+elisp\n"
	brokenSampleTwo   string = "INDEX\n"
	brokenSampleThree string = "QUERY|cloog|"
)

func TestSampleFullMessages(t *testing.T) {
	testMessage := func(sample string, message *message.Message, failureMessage string) {
		parsedMessage, err := parser.ParseMessage(sample)
		if err != nil {
			t.Errorf("Parsing of message %v failed because: %v", sample, err.Error())
		}
		if !message.Equals(parsedMessage) {
			message := fmt.Sprintf("%v;Expected: %v,Got: %v", failureMessage, message.String(), parsedMessage.String())
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
		parser.ParseMessage(sampleOne)
		parser.ParseMessage(sampleTwo)
		parser.ParseMessage(sampleThree)
		parser.ParseMessage(sampleFour)
	}
}

func TestSampleBrokenMessages(t *testing.T) {
	testMessage := func(sample string, failureMessage string) {
		_, err := parser.ParseMessage(sample)
		if err == nil {
			t.Errorf("Parsing of message %v passed.", sample)
		}
	}

	testMessage(brokenSampleOne, "INDEX without ending bar succeeded.")
	testMessage(brokenSampleTwo, "INDEX without immediately following bar succeeded.")
	testMessage(brokenSampleThree, "No newline separator succeeded.")
}
