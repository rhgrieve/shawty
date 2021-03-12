package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/db"
	"main/urls"
	"net/http"
	"net/url"

	"github.com/bmizerany/pat"
)

var database *db.DB

// SERVER

func urlHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	record, err := database.Get(id)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, record.JSON())
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	record, err := database.Get(id)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		database.Transaction(func() {
			record.IncrementVisits()
		})

		http.Redirect(w, r, record.URL, 302)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	records := database.JSON()

	fmt.Fprint(w, records)
}

type CreateRequest struct {
	URL string
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var c CreateRequest

	err := json.NewDecoder(r.Body).Decode(&c)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }

	unescapedURL, err := url.QueryUnescape(c.URL)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }

  validURL, err := url.ParseRequestURI(unescapedURL)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }

	convertedUrl := urls.NewShortURL(validURL)

	database.Transaction(func() {
		database.Add(convertedUrl)
	})

	record, err := database.Get(convertedUrl.Short)
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

	fmt.Fprint(w, "https://" + r.Host + "/" + record.URL.Short)
}

func main() {
	database = db.NewDB("db.json")
	database.Load()

	m := pat.New()

	// GET
	m.Get("/", http.HandlerFunc(rootHandler))
	m.Get("/urls/:id", http.HandlerFunc(urlHandler))
	m.Get("/:id", http.HandlerFunc(redirectHandler))

	// POST
	m.Post("/", http.HandlerFunc(createHandler))

	http.Handle("/", m)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}