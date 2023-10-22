package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"text/template"
)

func HandleEditPage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html", "templates/partials/head.html", "templates/edit.html",
	))

	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.Context().Value("filename").(string)
		id := r.Context().Value("id").(int64)

		err := tmpl.ExecuteTemplate(w, "base", struct {
			Filename   string
			EncodedURL string
		}{
			Filename:   filename,
			EncodedURL: url.QueryEscape(fmt.Sprintf("/files/%d", id)),
		})

		if err != nil {
			http.Error(w, "Something went wrong executing template", http.StatusInternalServerError)
		}
	}
}
