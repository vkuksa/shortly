package rest

import (
	"errors"
	"log"
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
	code, message := resolveErrorCode(err), err.Error()

	if err := c.metrics.CollectHTTPError(r.Method, r.URL.Path, strconv.Itoa(code), message); err != nil {
		log.Printf("[metrics] error: collection failed: %s", err.Error())
	}

	log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
	c.writeJSONResponse(w, r, NewErrResponse(message), code)
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
