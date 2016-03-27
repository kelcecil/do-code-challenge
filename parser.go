package main

import (
	"bytes"
	"io"
	"strings"
)

func ParseMessage(rawMessage string) (*Message, error) {
	newMessage := NewEmptyMessage()
	var buf bytes.Buffer
	buf.WriteString(rawMessage)

	commandToken, err := buf.ReadString(byte('|'))
	if err != nil {
		return nil, err
	}
	commandToken = strings.TrimRight(commandToken, "|")
	newMessage.Command = commandToken

	packageToken, err := buf.ReadString(byte('|'))
	if err != nil {
		return nil, err
	}
	packageToken = strings.TrimRight(packageToken, "|")
	newMessage.PackageName = packageToken

	var depBuf bytes.Buffer
	dependencyToken, err := buf.ReadString(byte('\n'))
	if err != nil {
		return nil, err
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
