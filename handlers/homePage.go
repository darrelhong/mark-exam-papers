package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/darrelhong/mark-exam-papers/utils"
)

func HandleHomePage() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html", "templates/partials/head.html", "templates/home.html",
	))

	return func(w http.ResponseWriter, r *http.Request) {

		files, err := utils.GetAllFiles()

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		err = tmpl.ExecuteTemplate(w, "base", struct {
			Files []utils.File
		}{
			Files: files,
		})

		if err != nil {
			http.Error(w, "Something went wrong executing template", http.StatusInternalServerError)
		}
	}
}
