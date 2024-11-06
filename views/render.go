package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// RenderTemplate carrega e renderiza o template com o layout base
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Caminhos para os templates
	basePath := filepath.Join("templates", "base.html")
	pagePath := filepath.Join("templates", tmpl)

	// Parseia o layout base e a página específica
	t, err := template.ParseFiles(basePath, pagePath)
	if err != nil {
		http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
		return
	}
	// Renderiza o template
	t.Execute(w, data)
}
