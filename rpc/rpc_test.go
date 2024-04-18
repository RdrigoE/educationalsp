package rpc_test

import (
	"educationalsp/rpc"
	"testing"
)

type EncondigExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncondigExample{Testing: true})
	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"

	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))

	contentLength := len(content)

	if err != nil {
		t.Fatal(err)
	}
	if contentLength != 15 {
		t.Fatalf("Expected: %d, Actual: %d", 15, contentLength)
	}
	if method != "hi" {
		t.Fatalf("Expected: 'hi', Got: %s", method)
	}
}
