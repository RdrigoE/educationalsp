package lsp

type CodeActionRequest struct {
	Request
	Params CodeActionParams `json:"params"`
}

type CodeActionParams struct {
	TextDocumentPositionParams
}

type CodeActionResponse struct {
	Response
	Result CodeActionResult `json:"result"`
}

type CodeActionResult struct {
	Contents string `json:"contents"`
}
