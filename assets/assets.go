package assets

import (
	"embed"
	"html/template"
	"io"
)

//go:embed *
var StaticFiles embed.FS

var (
	index = template.Must(template.ParseFS(StaticFiles, "index.html"))
)

func RenderIndex(w io.Writer, data any) error {
	return index.Execute(w, data)
}
