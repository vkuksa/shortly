package assets

import (
	"embed"
	"html/template"
	"io"
)

//go:embed all:*.*
var files embed.FS

var (
	index = load("index.html")
)

type DashboardParams struct {
	Title   string
	Message string
}

func RenderIndex(w io.Writer, data any) error {
	return index.Execute(w, data)
}

func load(file string) *template.Template {
	return template.Must(template.New("index.html").ParseFS(files, file))
}
