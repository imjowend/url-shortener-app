package controllers

import (
	"database/sql"
	"github.com/imjowend/url-shortener-app/internal/db"
	"github.com/imjowend/url-shortener-app/internal/url"
	"html/template"
	"net/http"
	"strings"
)

func Shorten(lite *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		originalURL := r.FormValue("url")
		if originalURL == "" {
			http.Error(w, "URL not provided", http.StatusBadRequest)
		}

		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			originalURL = "http://" + originalURL
		}

		shortURL := url.Shorten(originalURL)

		if err := db.StoreURL(lite, shortURL, originalURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]string{
			"ShortURL": shortURL,
		}

		t, err := template.ParseFiles("internal/views/shorten.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func Proxy(lite *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.URL.Path[1:]
		if shortURL == "" {
			http.Error(w, "URL not provided", http.StatusBadRequest)
			return
		}
		originalURL, err := db.GetOriginalURL(lite, shortURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
	}
}
