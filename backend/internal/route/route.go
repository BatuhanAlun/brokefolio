package route

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var tpl *template.Template

func InitTemplates() {
	templateDir := "./templates"
	templatePattern := filepath.Join(templateDir, "*.html")

	tpl = template.Must(template.ParseGlob(templatePattern))
	log.Println("Templates Loaded successfully")
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := tpl.ExecuteTemplate(w, tmpl, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error on rendering Template %s %v", tmpl, err)
	}

}
