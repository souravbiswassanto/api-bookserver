package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

//func getBooks

type Author struct {
	//name string `json:"name,omitempty"`
	//bookCount string `json:"book_count"`
	age string `json:"age"`
	//books     []string `json:"books"`
}

type Book struct {
	name    string   `json:"book_name,omitempty"`
	authors []Author `json:"authors"`
	isbn    string   `json:"isbn,omitempty"`
	genre   string   `json:"genre"`
}

type bookDB map[string]Book
type authorDB map[string]Author

var bookList bookDB
var authorList authorDB

type test struct {
	tmp string `json:"tmp"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	v := test{tmp: "hello"}
	fmt.Println(v.tmp)
	w.Write([]byte("Pong"))
	json.NewEncoder(w).Encode(v)

	return
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//body, _ := io.ReadAll(r.Body)
	//fmt.Println("Request body:", string(body))
	//fmt.Println(book.age)
	//json.NewEncoder(w).Encode(book)
	var book Author

	json.NewDecoder(r.Body).Decode(&book)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("appscode"))
	//json.NewEncoder(w).Encode(v)
	fmt.Println("========###end###=====")
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/book", Ping)
	r.Get("/books", getBooks)
	/*
		r.Group(func(r chi.Router) {
			r.Route("/books", func(r chi.Router) {
				r.Get("/", getBooks)
			})
		})
	*/
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}
