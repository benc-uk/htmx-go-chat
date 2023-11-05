package main

// Implement a simple renderer for the echo framework

import (
	"bytes"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// HTMLRenderer is a custom renderer for echo
type HTMLRenderer struct {
	templates *template.Template
}

// NewHTMLRenderer creates a new renderer
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

// Render renders a template document
func (r *HTMLRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

// RenderToString renders a template document to a string
func (r *HTMLRenderer) RenderToString(name string, data interface{}) (string, error) {
	var buf []byte
	w := bytes.NewBuffer(buf)
	err := r.Render(w, name, data, nil)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}
