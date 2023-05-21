package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vkuksa/shortly/assets"
)

func (s *Server) registerLinkRoutes(router chi.Router) {
	router.Get("/", s.handleRoot)
	router.Post("/", s.handleStoreLink)
	router.Get("/{uuid}", s.handleRedirrectLink)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if err := assets.RenderIndex(w, nil); err != nil {
		s.handleError(w, r, err)
	}
}

func (s *Server) handleStoreLink(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	link, err := s.LinkService.GenerateShortenedLink(r.Context(), url)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	l := struct {
		Shortened string
		Original  string
	}{
		Shortened: s.url() + "/" + link.UUID,
		Original:  link.URL,
	}

	if err := assets.RenderIndex(w, l); err != nil {
		s.handleError(w, r, err)
		return
	}
}

func (s *Server) handleRedirrectLink(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	link, err := s.LinkService.GetOriginalLink(r.Context(), uuid)
	if err != nil {
		s.handleError(w, r, err)
		return
	}

	if err = s.LinkService.AddHit(r.Context(), uuid); err != nil {
		s.logError(r, err)
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}
