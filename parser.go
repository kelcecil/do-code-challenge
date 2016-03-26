package main

import "bytes"

func ParseMessage(rawMessage string) (*Message, error) {
	newMessage := NewEmptyMessage()
	var buf bytes.Buffer

	for i := 0; i < len(rawMessage); i++ {
		buf.WriteByte(rawMessage[i])
		byteBar := byte('|')
		if rawMessage[i] == byteBar {
			token, err := buf.ReadBytes(byteBar)
			if err != nil {
				return nil, err
			}

			token = bytes.TrimSuffix(token, []byte("|"))
			if newMessage.Command == "" {
				newMessage.Command = string(token)
			} else if newMessage.PackageName == "" {
				newMessage.PackageName = string(token)
			}
		} else if rawMessage[i] == byte(',') {
			token, err := buf.ReadBytes(byte(','))
			if err != nil {
				return nil, err
			}
			token = bytes.TrimSuffix(token, []byte(","))
			newMessage.PackageDependencies = append(newMessage.PackageDependencies, string(token))
		} else if rawMessage[i] == byte('\n') {
			token, err := buf.ReadBytes(byte('\n'))
			if err != nil {
				return nil, err
			}
			token = bytes.TrimSpace(bytes.TrimSuffix(token, []byte("\n")))
			if len(token) != 0 {
				newMessage.PackageDependencies = append(newMessage.PackageDependencies, string(token))
			}
			break
		}
	}
	return newMessage, nil
}
