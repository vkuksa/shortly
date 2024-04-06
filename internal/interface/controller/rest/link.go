package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/vkuksa/shortly/internal/link"

	"github.com/go-chi/chi/v5"
)

type ErrorHandler interface {
	HandleRESTError(w http.ResponseWriter, r *http.Request, err error)
}

type LinkController struct {
	uc         *link.UseCase
	errhandler ErrorHandler
}

func NewLinkController(uc *link.UseCase, eh ErrorHandler) *LinkController {
	return &LinkController{uc: uc, errhandler: eh}
}

func (c *LinkController) Register(router chi.Router) {
	router.Post("/links", c.handleStore)
	router.Get("/links/{uuid}", c.handleRetrieve)
	router.Get("/{uuid}", c.handleRedirrect)
}

type LinkRequest struct {
	URL string `json:"url"`
}

func (c *LinkController) handleStore(w http.ResponseWriter, r *http.Request) {
	var req LinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.errhandler.HandleRESTError(w, r, link.ErrBadInput)
		return
	}

	link, err := c.uc.Shorten(r.Context(), req.URL)
	if err != nil {
		c.errhandler.HandleRESTError(w, r, err)
		return
	}

	c.writeJSONResponse(w, link, http.StatusOK)
}

func (c *LinkController) handleRedirrect(w http.ResponseWriter, r *http.Request) { //TODO: should it be here? make only web-based redirrects
	uuid := chi.URLParam(r, "uuid")
	link, err := c.uc.Retrieve(r.Context(), uuid)
	if err != nil {
		c.errhandler.HandleRESTError(w, r, err)
		return
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}

func (c *LinkController) handleRetrieve(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	link, err := c.uc.Retrieve(r.Context(), uuid)
	if err != nil {
		c.errhandler.HandleRESTError(w, r, err)
		return
	}

	c.writeJSONResponse(w, link, http.StatusOK)
}

func (c *LinkController) writeJSONResponse(w http.ResponseWriter, obj any, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		slog.Error("writing response failed", slog.Any("component", "rest"), slog.Any("error", err))
		return
	}
}
