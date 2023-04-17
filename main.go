package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

//func getBooks

type Author struct {
	Name      string   `json:"name,omitempty"`
	BookCount string   `json:"book_count"`
	Age       string   `json:"age"`
	Books     []string `json:"books"`
}

type Book struct {
	Name    string   `json:"book_name,omitempty"`
	Authors []Author `json:"authors"`
	ISBN    string   `json:"isbn,omitempty"`
	Genre   string   `json:"genre"`
	Pub     string   `json:"publisher"`
}

type BookDB map[string]Book
type AuthorDB map[string]Author

var BookList BookDB
var AuthorList AuthorDB

func init() {
	author1 := Author{
		Name:      "temp author 1",
		BookCount: "5",
		Age:       "45",
		Books:     []string{"ISBN 1", "ISBN 2"},
	}
	author2 := Author{
		Name:      "temp author 2",
		BookCount: "5",
		Age:       "45",
		Books:     []string{"ISBN 3", "ISBN 4"},
	}

	data1 := Book{
		Name: "temp book 1",
		Authors: []Author{
			author1,
			author2,
		},
		ISBN:  "ISBN 1",
		Genre: "Fiction",
		Pub:   "Demo",
	}

	data2 := Book{
		Name: "temp book 2",
		Authors: []Author{
			author1,
		},
		ISBN:  "ISBN 2",
		Genre: "Fiction",
		Pub:   "Demo",
	}

	data3 := Book{
		Name: "temp book 3",
		Authors: []Author{
			author2,
		},
		ISBN:  "ISBN 3",
		Genre: "Fiction",
		Pub:   "Demo",
	}
	BookList = make(BookDB)
	BookList[data1.ISBN] = data1
	BookList[data2.ISBN] = data2
	BookList[data3.ISBN] = data3

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//var book Book
	//json.NewDecoder(r.Body).Decode(&book)
	//json.NewEncoder(w).Encode(book)

	err := json.NewEncoder(w).Encode(BookList)

}

func main() {
	init()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Group(func(r chi.Router) {
		r.Route("/books", func(r chi.Router) {
			r.Get("/", getBooks)
		})
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}
