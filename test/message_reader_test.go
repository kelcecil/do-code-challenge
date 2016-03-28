package test

import (
	"testing"

	"../message"
	"../parser"
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

func TestMessageReaderIndexDependencies(t *testing.T) {
	msg := message.NewMessage("INDEX", "golang", []string{"glide", "yaml"})
	testMessageReader(t, "INDEX|golang|glide,yaml\n", msg)
}

func TestMessageReaderIndex(t *testing.T) {
	msg := message.NewMessage("INDEX", "golang", []string{})
	testMessageReader(t, "INDEX|golang|\n", msg)
}

func TestMessageReaderRemove(t *testing.T) {
	msg := message.NewMessage("REMOVE", "golang", []string{})
	testMessageReader(t, "REMOVE|golang|\n", msg)
}

func TestMessageReaderQuery(t *testing.T) {
	msg := message.NewMessage("QUERY", "golang", []string{})
	testMessageReader(t, "QUERY|golang|\n", msg)
}

func testMessageReader(t *testing.T, query string, expected *message.Message) {
	call := func(b []byte) (n int, err error) {
		result := []byte(query)
		for i := range result {
			b[i] = result[i]
		}
		return len(result), nil
	}

	mcr := NewMockConnectionReader(call)
	reader := parser.NewMessageReader(mcr)

	msg, err := reader.Read()
	if err != nil {
		t.Errorf("Reader returned an error; Error: %v", err)
	}

	if !msg.Equals(expected) {
		t.Errorf("Messages did not match; Wanted: %v, Got: %v", expected.String(), msg.String())
	}
}
