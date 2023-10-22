package handlers

import (
	"net/http"
	"text/template"
)

func HandleUploadPage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html", "templates/partials/head.html", "templates/upload.html",
	))

	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "base", nil)

		if err != nil {
			http.Error(w, "Something went wrong executing template", http.StatusInternalServerError)
		}
	}
}
