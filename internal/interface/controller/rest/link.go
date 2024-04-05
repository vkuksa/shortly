package rest

import (
	"encoding/json"
	"net/http"

	"github.com/vkuksa/shortly/internal/link"

	"github.com/go-chi/chi/v5"
)

type MetricsCollector interface {
	CollectHttpError(method, path string, labels ...string) error
}

type LinkController struct {
	uc      *link.UseCase
	metrics MetricsCollector
}

func NewLinkController(uc *link.UseCase, mc MetricsCollector) *LinkController {
	return &LinkController{uc: uc, metrics: mc}
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
		c.handleError(w, r, link.ErrBadInput)
		return
	}

	link, err := c.uc.Shorten(r.Context(), req.URL)
	if err != nil {
		c.handleError(w, r, err)
		return
	}

	c.writeJsonResponse(w, r, link, http.StatusOK)
}

func (c *LinkController) handleRedirrect(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	link, err := c.uc.Retrieve(r.Context(), uuid)
	if err != nil {
		c.handleError(w, r, err)
		return
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}

func (c *LinkController) handleRetrieve(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	link, err := c.uc.Retrieve(r.Context(), uuid)
	if err != nil {
		c.handleError(w, r, err)
		return
	}

	c.writeJsonResponse(w, r, link, http.StatusOK)
}

func (c *LinkController) writeJsonResponse(w http.ResponseWriter, r *http.Request, obj any, status int) {
	data, err := json.Marshal(obj)
	if err != nil {
		c.handleError(w, r, err)
		return
	}

	w.WriteHeader(status)
	w.Write(data)
}
