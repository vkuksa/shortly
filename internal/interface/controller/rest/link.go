package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/vkuksa/shortly/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type LinkController struct {
	uc usecase.UseCase
}

func NewLinkController(uc usecase.UseCase) LinkController {
	return LinkController{uc: uc}
}

func (c LinkController) Register(router chi.Router) {
	router.Post("/links", c.handleStore)
	router.Get("/links/{shortened}", c.handleRetrieve)
	router.Get("/{shortened}", c.handleRedirrect)
}

type LinkRequest struct {
	URL string `json:"url"`
}

func (c LinkController) handleStore(w http.ResponseWriter, r *http.Request) {
	var req LinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// c.errhandler.HandleRESTError(w, r, link.ErrBadInput)
		return
	}

	link, err := c.uc.ShortenUrl(r.Context(), usecase.NewShortenUrlInput(req.URL))
	if err != nil {
		// c.errhandler.HandleRESTError(w, r, err)
		return
	}

	c.writeJSONResponse(w, link, http.StatusCreated)
}

func (c LinkController) handleRedirrect(w http.ResponseWriter, r *http.Request) { //TODO: should it be here? make only web-based redirrects
	shortened := chi.URLParam(r, "shortened")
	link, err := c.uc.GetOriginal(r.Context(), usecase.NewGetOriginalInput(shortened))
	if err != nil {
		// c.errhandler.HandleRESTError(w, r, err)
		return
	}

	http.Redirect(w, r, link.Original, http.StatusFound)
}

func (c *LinkController) handleRetrieve(w http.ResponseWriter, r *http.Request) {
	shortened := chi.URLParam(r, "shortened")
	link, err := c.uc.GetOriginalWithoutHit(r.Context(), usecase.NewGetOriginalInput(shortened))
	if err != nil {
		// c.errhandler.HandleRESTError(w, r, err)
		return
	}

	c.writeJSONResponse(w, link, http.StatusOK)
}

func (c LinkController) writeJSONResponse(w http.ResponseWriter, obj any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		slog.Error("writing response failed", slog.Any("component", "rest"), slog.Any("error", err))
		return
	}
}
