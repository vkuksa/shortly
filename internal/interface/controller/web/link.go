package web

import (
	"log/slog"
	"net/http"

	"github.com/vkuksa/shortly/assets"

	"github.com/go-chi/chi/v5"
)

type LinkController struct {
}

func NewLinkController() LinkController {
	return LinkController{}
}

func (c LinkController) Register(router chi.Router) {
	fs := http.FS(assets.StaticFiles)
	router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(fs)))

	router.Get("/", c.handleHome)
}

func (c LinkController) handleHome(w http.ResponseWriter, r *http.Request) {
	data := struct{}{}
	if err := assets.RenderIndex(w, data); err != nil {
		slog.Error("Template rendering error", slog.Any("error", err))
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
