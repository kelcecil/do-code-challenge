package parser

import (
	"bytes"
	"errors"
	"github.com/kelcecil/do-code-challenge/message"
	"io"
	"strings"
)

var (
	PARSE_FAILED error = errors.New("Parsing of message failed.")
)

type MessageReader struct {
	reader io.Reader
	buf    bytes.Buffer
}

func NewMessageReader(rdr io.Reader) MessageReader {
	return MessageReader{
		reader: rdr,
	}
}

func (rdr *MessageReader) Read() (*message.Message, error) {
	readBytes := make([]byte, 4096)
	for {
		n, err := rdr.reader.Read(readBytes)
		if err == io.EOF {
			return nil, err
		}
		rdr.buf.Write(readBytes[:n])
		if bytes.Contains(readBytes, []byte("\n")) {
			break
		}
	}

	line, _ := rdr.buf.ReadString(byte('\n'))
	message, err := ParseMessage(line)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func ParseMessage(rawMessage string) (*message.Message, error) {
	newMessage := message.NewEmptyMessage()
	var buf bytes.Buffer
	buf.WriteString(rawMessage)

	commandToken, err := buf.ReadString(byte('|'))
	if err != nil {
		return nil, PARSE_FAILED
	}
	commandToken = strings.TrimRight(commandToken, "|")
	newMessage.Command = commandToken

	packageToken, err := buf.ReadString(byte('|'))
	if err != nil {
		return nil, PARSE_FAILED
	}
	packageToken = strings.TrimRight(packageToken, "|")
	newMessage.PackageName = packageToken

	var depBuf bytes.Buffer
	dependencyToken, err := buf.ReadString(byte('\n'))
	if err != nil {
		return nil, PARSE_FAILED
	}
	depBuf.WriteString(dependencyToken)
	for {
		token, err := depBuf.ReadString(byte(','))
		token = strings.TrimSpace(strings.TrimRight(token, ","))
		if len(token) != 0 {
			newMessage.PackageDependencies = append(newMessage.PackageDependencies, token)
		}
		if err == io.EOF {
			break
		}
	}
	return newMessage, nil
}
