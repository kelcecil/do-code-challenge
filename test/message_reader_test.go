package test

import (
	"github.com/kelcecil/do-code-challenge/message"
	"github.com/kelcecil/do-code-challenge/parser"
	"testing"
)

type MockConnectionReader struct {
	testFunc func([]byte) (int, error)
}

type MockReaderFunction func(b []byte) (n int, err error)

func NewMockConnectionReader(reader MockReaderFunction) MockConnectionReader {
	return func() MockConnectionReader {
		return MockConnectionReader{
			testFunc: reader,
		}
	}()
}

func (mcr MockConnectionReader) Read(b []byte) (n int, err error) {
	return mcr.testFunc(b)
}

func TestMessageReaderIndex(t *testing.T) {
	index := func(b []byte) (n int, err error) {
		result := []byte("INDEX|golang|\n")
		for i := range result {
			b[i] = result[i]
		}
		return len(result), nil
	}

	mcr := NewMockConnectionReader(index)
	reader := parser.NewMessageReader(mcr)

	msg, err := reader.Read()
	if err != nil {
		t.Errorf("Reader returned an error; Error: %v", err)
	}
	expectedMessage := message.NewMessage("INDEX", "golang", []string{})

	if !msg.Equals(expectedMessage) {
		t.Error("Messages did not match; Wanted: %v, Got: %v", expectedMessage.String(), msg.String())
	}
}
