package http

import (
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"
)

// TODO: remove this function and usages and think of getting proper template path
func indexPath() (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	moduleRoot := strings.TrimSpace(string(output))

	// Construct the template file path relative to the current working directory
	templatePath := filepath.Join(moduleRoot, "web", "index.html")

	return templatePath, nil
}

func (s *Server) registerLinkRoutes() {
	s.router.Get("/", s.handleRoot)
	s.router.Post("/", s.handleStoreLink)
	s.router.Get("/{uuid}", s.handleRedirrectLink)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	path, err := indexPath()
	if err != nil {
		s.HandleError(w, r, err)
		return
	}

	tmpl := template.Must(template.ParseFiles(path))

	if err := tmpl.Execute(w, nil); err != nil {
		s.HandleError(w, r, err)
		return
	}
}

func (s *Server) handleStoreLink(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	link, err := s.LinkService.GenerateShortenedLink(r.Context(), url)
	if err != nil {
		s.HandleError(w, r, err)
		return
	}

	l := struct {
		Shortened string
		Original  string
	}{
		Shortened: s.URL() + "/" + link.UUID,
		Original:  link.URL,
	}

	path, err := indexPath()
	if err != nil {
		s.HandleError(w, r, err)
		return
	}

	tmpl := template.Must(template.ParseFiles(path))

	if err := tmpl.Execute(w, l); err != nil {
		s.HandleError(w, r, err)
		return
	}
}

func (s *Server) handleRedirrectLink(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	link, err := s.LinkService.GetOriginalLink(r.Context(), uuid)
	if err != nil {
		s.HandleError(w, r, err)
		return
	}

	if err = s.LinkService.AddHit(r.Context(), uuid); err != nil {
		s.LogError(r, err)
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}
