package stub

import "net/http"

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (c *ErrorHandler) HandleRESTError(_ http.ResponseWriter, _ *http.Request, _ error) {}
