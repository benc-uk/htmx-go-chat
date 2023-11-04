package main

// Implement a simple renderer for the echo framework

import (
	"bytes"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// Renderer is a custom renderer for echo
type Renderer struct {
	templates *template.Template
}

// HTMLRenderer creates a new renderer
func HTMLRenderer() *Renderer {
	return &Renderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

// Render renders a template document
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

func (r *Renderer) RenderToString(name string, data interface{}) (string, error) {
	var buf []byte
	w := bytes.NewBuffer(buf)
	err := r.Render(w, name, data, nil)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}
