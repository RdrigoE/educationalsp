package main

import (
	"bufio"
	"educationalsp/analysis"
	"educationalsp/lsp"
	"educationalsp/rpc"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("hi")
	logger := getLogger("/home/reusebio/projets/educationalsp/log.txt")
	logger.Println("Hey, I started!")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("ERROR: %s", err)
			continue
		}
		handleMessage(logger, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, state analysis.State, method string, contents []byte) {
	logger.Println("Received msg with method: ", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Not possible to parse. CONTENT: %s", err)
		}

		logger.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)
		msg := lsp.NewInitilizeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)
		writer := os.Stdout
		writer.Write([]byte(reply))
		logger.Print("Reply to initialize")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didOpen: %s", err)
		}

		logger.Printf("Open: %s", request.Params.TextDocument.URI)

		// Sync the state
		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didChange: %s", err)
			return
		}

		logger.Printf("Changed: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
		}

		// Sync the state
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Give me a good file")
	}
	return log.New(logfile, "[educaationalsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
