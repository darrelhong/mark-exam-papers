package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/darrelhong/mark-exam-papers/utils"
)

func HandleGetFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.Context().Value("filename").(string)
		http.ServeFile(w, r, "storage/"+filename)
	}
}

func HandleDeleteFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.Context().Value("filename").(string)
		id := r.Context().Value("id").(int64)
		err := os.Remove("storage/" + filename)

		if err != nil {
			http.Error(w, "Something went wrong deleting the file", http.StatusInternalServerError)
		}

		err = utils.DeleteFile(id)
		if err != nil {
			http.Error(w, "Something went wrong deleting the file", http.StatusInternalServerError)
		}
	}
}

func HandleUploadFile() http.HandlerFunc {
	err := os.MkdirAll("storage", 0755)

	if err != nil {
		fmt.Println("Something went wrong creating the storage directory")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body to retrieve the files
		r.ParseMultipartForm(32 << 20) // 32 MB is the maximum memory used to parse the request body
		files := r.MultipartForm.File["files"]

		// store the files in the local storage
		for _, file := range files {
			fmt.Println(file.Filename)
			newFilename := utils.SanitiseFilename(file.Filename)
			fmt.Println(newFilename)

			f, err := file.Open()
			if err != nil {
				fmt.Fprintf(w, "Something went wrong opening the file")
				return
			}
			defer f.Close()

			newFile, err := os.Create("storage/" + newFilename)
			if err != nil {
				fmt.Fprintf(w, "Something went wrong creating the file")
				return
			}
			defer newFile.Close()

			_, err = io.Copy(newFile, f)
			if err != nil {
				fmt.Fprintf(w, "Something went wrong copying the file")
				return
			}

			// add the file to the database
			utils.DB.Exec("INSERT INTO files (filename) VALUES (?)", newFilename)
		}

		// write hx-redirect header to redirect to the home page
		w.Header().Set("hx-redirect", "/")
	}
}
