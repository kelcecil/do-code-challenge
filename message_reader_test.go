package main

import (
	"testing"
)

type MockConnectionReader struct {
	NextResult func([]byte) (int, error)
}

func NewMockConnectionReader() MockConnectionReader {
	return MockConnectionReader{
		NextResult: func(b []byte) (n int, err error) {
			result := []byte("INDEX|golang|\n")
			for i := range result {
				b[i] = result[i]
			}
			return len(result), nil
		},
	}
}

func (mcr MockConnectionReader) Read(b []byte) (n int, err error) {
	return mcr.NextResult(b)
}

func TestMessageReader(t *testing.T) {
	mcr := NewMockConnectionReader()
	reader := NewMessageReader(mcr)

	message, err := reader.Read()
	if err != nil {
		t.Errorf("Reader returned an error; Error: %v", err)
	}
	expectedMessage := NewMessage("INDEX", "golang", []string{})

	if !message.Equals(expectedMessage) {
		t.Error("Messages did not match; Wanted: %v, Got: %v", expectedMessage.String(), message.String())
	}
}
