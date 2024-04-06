package errhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/vkuksa/shortly/internal/link"
)

type MetricsCollector interface {
	CollectHTTPError(method, path string, labels ...string)
}

type ErrorHandler struct {
	metrics MetricsCollector
}

func NewErrorHandler(metrics MetricsCollector) *ErrorHandler {
	return &ErrorHandler{metrics: metrics}
}

func (h *ErrorHandler) HandleRESTError(w http.ResponseWriter, r *http.Request, err error) {
	h.handleError("rest", w, r, err)
}

func (h *ErrorHandler) HandleGraphQLError(w http.ResponseWriter, r *http.Request, err error) {
	h.handleError("gql", w, r, err)
}

type errResponse struct {
	Errors string `json:"errors"`
}

func newErrResponse(m string) errResponse {
	return errResponse{Errors: m}
}

func (c *ErrorHandler) handleError(component string, w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(fmt.Sprintf("%s\t%s", r.Method, r.URL.Path), slog.Any("component", component), slog.Any("error", err))

	code, message := resolveErrorCode(err), err.Error()
	c.writeJSONResponse(w, newErrResponse(message), code)

	c.metrics.CollectHTTPError(r.Method, r.URL.Path, strconv.Itoa(code), message)
}

func (c *ErrorHandler) writeJSONResponse(w http.ResponseWriter, obj any, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		slog.Error("writing response failed", slog.Any("component", "errhandler"), slog.Any("error", err))
		return
	}
}
func resolveErrorCode(err error) int {
	switch {
	case errors.Is(err, link.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, link.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, link.ErrBadInput):
		return http.StatusBadRequest
	case errors.Is(err, link.ErrInternal):
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
