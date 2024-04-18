package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("Did not find separator")
	}
	// Content-Lengtg: <number>
	contentLenghtBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLenghtBytes))
	if err != nil {
		return "", nil, err
	}
	// TODO: We will get to this
	var baseMessage BaseMessage
	contentB := content[:contentLength]
	if err := json.Unmarshal(contentB, &baseMessage); err != nil {
		return "", nil, err
	}
	return baseMessage.Method, contentB, nil
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
func Split(data []byte, _ bool) (advance int, token []byte, err error) {

	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil // Not a problem, we just need more data
	}
	// Content-Lengtg: <number>
	contentLenghtBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLenghtBytes))
	if err != nil {
		return 0, nil, err
	}
	if contentLength < len(content) {
		return 0, nil, nil // Not a problem, we just need more data
	}
	totalLenght := len(header) + 4 + contentLength
	// TODO: We will get to this
	return totalLenght, data[:totalLenght], nil
}
