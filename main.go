package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/darrelhong/mark-exam-papers/handlers"
	"github.com/darrelhong/mark-exam-papers/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	utils.InitDb()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handlers.HandleHomePage())

	r.Route("/files", func(r chi.Router) {
		r.Get("/", handlers.HandleUploadPage())
		r.Post("/", handlers.HandleUploadFile())

		r.Route("/{id}", func(r chi.Router) {
			r.Use(fileContext)
			r.Get("/", handlers.HandleGetFile())
			r.Delete("/", handlers.HandleDeleteFile())
			r.Get("/edit", handlers.HandleEditPage())
		})
	})

	fileServer(r, "/static", http.Dir("static"))
	http.ListenAndServe("localhost:8080", r)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

// extract the filename from the URL and add it to the context
// error if no filename is provided
func fileContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		if idStr == "" {
			http.Error(w, "No id provided", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		//get the filename from the id
		filename, err := utils.GetFilename(id)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		ctx := context.WithValue(r.Context(), "filename", filename)
		ctx = context.WithValue(ctx, "id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
