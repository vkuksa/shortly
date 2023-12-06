package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/vkuksa/shortly/assets"
	"github.com/vkuksa/shortly/internal/domain"
	"github.com/vkuksa/shortly/internal/infrastructure/metrics"
	"github.com/vkuksa/shortly/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type LinkUseCase interface {
	GenerateShortenedLink(ctx context.Context, url string) (*domain.Link, error)
	GetOriginalLink(ctx context.Context, url string) (*domain.Link, error)
}

type LinkController struct {
	uc LinkUseCase
}

func NewLinkController(uc LinkUseCase) *LinkController {
	return &LinkController{uc: uc}
}

func (c *LinkController) Register(router chi.Router) {
	router.Get("/", c.handleRoot)
	router.Post("/", c.handleStoreLink)
	router.Get("/{uuid}", c.handleRedirrectLink)
}

func (c *LinkController) handleRoot(w http.ResponseWriter, r *http.Request) {
	if err := assets.RenderIndex(w, nil); err != nil {
		c.handleError(w, r, err)
	}
}

func (c *LinkController) handleStoreLink(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	link, err := c.uc.GenerateShortenedLink(r.Context(), url)
	if err != nil {
		c.handleError(w, r, err)
		return
	}

	l := struct {
		Shortened string
		Original  string
	}{
		Shortened: "http://localhost:8080/" + link.UUID, //TODO: dynamic url here
		Original:  link.URL,
	}

	if err := assets.RenderIndex(w, l); err != nil {
		c.handleError(w, r, err)
		return
	}
}

func (c *LinkController) handleRedirrectLink(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	log.Print("uuid")
	log.Print(uuid)

	link, err := c.uc.GetOriginalLink(r.Context(), uuid)
	if err != nil {
		c.handleError(w, r, err)
		return
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}

func (c *LinkController) handleError(w http.ResponseWriter, r *http.Request, err error) {
	code, message := usecase.ErrorCode(err), usecase.ErrorMessage(err)

	metrics.ErrorCount.WithLabelValues(r.Method, r.URL.Path, code, message).Inc()

	log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
	http.Error(w, message, errorStatusCode(code))
}

var codes = map[string]int{
	usecase.ErrConflict: http.StatusConflict,
	usecase.ErrInvalid:  http.StatusBadRequest,
	usecase.ErrNotFound: http.StatusNotFound,
	usecase.ErrInternal: http.StatusInternalServerError,
}

func errorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}

	return http.StatusInternalServerError
}
