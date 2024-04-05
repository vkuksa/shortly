package rest

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/vkuksa/shortly/internal/link"
)

type ErrResponse struct {
	Message string `json:"message"`
}

func NewErrResponse(m string) ErrResponse {
	return ErrResponse{Message: m}
}

func (c *LinkController) handleError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(fmt.Sprintf("%s\t%s", r.Method, r.URL.Path), slog.Any("component", "rest"), slog.Any("error", err))

	code, message := resolveErrorCode(err), err.Error()
	c.writeJSONResponse(w, r, NewErrResponse(message), code)

	c.metrics.CollectHTTPError(r.Method, r.URL.Path, strconv.Itoa(code), message)
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
