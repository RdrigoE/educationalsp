package main

import (
	"bufio"
	"educationalsp/analysis"
	"educationalsp/lsp"
	"educationalsp/rpc"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fmt.Println("hi")
	logger := getLogger("/home/reusebio/projects/educationalsp/log.txt")
	logger.Println("Hey, I started!")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	writer := os.Stdout
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("ERROR: %s", err)
			continue
		}
		logger.Print(method)
		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
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
		writeResponse(writer, msg)
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
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover: %s", err)
		}
		// Create a response and write back
		msg := state.Hover(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		writeResponse(writer, msg)

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/definition: %s", err)
		}
		// Create a response and write back
		msg := state.Definition(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		writeResponse(writer, msg)

	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/definition: %s", err)
		}
		// Create a response and write back
		msg := state.TextDocumentCompletion(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		writeResponse(writer, msg)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Give me a good file")
	}
	return log.New(logfile, "[educaationalsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
