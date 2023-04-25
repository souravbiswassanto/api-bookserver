package apiHandler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/souravbiswassanto/api-bookserver/authHandler"
	"github.com/souravbiswassanto/api-bookserver/dataHandler"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetBooks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(dataHandler.BookList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func BookGeneralized(w http.ResponseWriter, _ *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	var readString []string
	for _, str := range dataHandler.BookList {
		readString = append(readString, str.Name)
	}
	resp := strings.Join(readString, "\n")
	if resp == "" {
		http.Error(w, "No Books found", http.StatusNotFound)
		return
	}
	_, err := w.Write([]byte(resp))

	if err != nil {
		http.Error(w, "Cannot Write Response", http.StatusInternalServerError)
		return
	}
}

func NewBook(w http.ResponseWriter, r *http.Request) {
	var book dataHandler.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Cannot Decode the data", http.StatusBadRequest)
		return
	}
	_, okay := dataHandler.BookList[book.ISBN]
	if len(book.Name) == 0 || len(book.ISBN) == 0 || len(book.Authors) == 0 || okay == true {
		http.Error(w, "Invalid Data Entry", http.StatusBadRequest)
		return
	}

	flag := false
	for _, data := range book.Authors {
		if len(data.Name) == 0 {
			flag = true
			break
		}
	}
	if flag == true {
		http.Error(w, "Author name can't be empty", http.StatusBadRequest)
		return
	}

	for _, author := range book.Authors {
		name := author.Name
		_, ok := dataHandler.AuthorList[dataHandler.SmStr(name)]
		var ab dataHandler.AuthorBooks
		if ok == false {
			ab.Author = author
			ab.Books = append(ab.Books, book.ISBN)
			dataHandler.AuthorList[dataHandler.SmStr(name)] = ab
			continue
		}
		ab = dataHandler.AuthorList[dataHandler.SmStr(name)]
		ab.Books = append(ab.Books, book.ISBN)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Data added successfully"))
	if err != nil {
		http.Error(w, "Can not Write Data", http.StatusInternalServerError)
		return
	}

	dataHandler.BookList[book.ISBN] = book
	//GetBooks(w, r)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	// receiving isbn of the book to be deleted
	var ISBN string
	ISBN = chi.URLParam(r, "ISBN")
	if len(ISBN) == 0 {
		http.Error(w, "ISBN is wrong", http.StatusBadRequest)
		return
	}
	_, ok := dataHandler.BookList[ISBN]
	if ok == false {
		http.Error(w, "ISBN is wrong", http.StatusBadRequest)
		return
	}

	delete(dataHandler.BookList, ISBN)

	for _, data := range dataHandler.AuthorList {
		var ab dataHandler.AuthorBooks
		ab = data
		for i, val := range ab.Books {
			if val == ISBN {
				ab.Books = append(ab.Books[:i], ab.Books[i+1:]...)
				ab.Books = ab.Books[:len(ab.Books)-1]
				// memory safety
				break
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Data added successfully"))
	if err != nil {
		http.Error(w, "Can not Write Data", http.StatusInternalServerError)
		return
	}
	fmt.Println("hello")

}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	// receiving isbn of the book to be deleted
	var ISBN string
	ISBN = chi.URLParam(r, "ISBN")
	fmt.Println(ISBN)
	if len(ISBN) == 0 {
		http.Error(w, "ISBN is wrong", http.StatusBadRequest)
		return
	}
	_, ok := dataHandler.BookList[ISBN]
	if ok == false {
		http.Error(w, "ISBN is wrong", http.StatusBadRequest)
		return
	}

	curBook := dataHandler.BookList[ISBN]
	var updBook dataHandler.Book

	err := json.NewDecoder(r.Body).Decode(&updBook)

	if err != nil {
		http.Error(w, "Cannot Decode the Given Data", http.StatusBadRequest)
		return
	}

	if curBook.ISBN != updBook.ISBN {
		http.Error(w, "Provided ISBN does not match", http.StatusBadRequest)
		return
	}

	// deleting old authors
	for _, author := range curBook.Authors {
		var ab dataHandler.AuthorBooks
		ab = dataHandler.AuthorList[dataHandler.SmStr(author.Name)]
		fmt.Println(ab)
		for i, val := range ab.Books {
			if val == ISBN {
				ab.Books = append(ab.Books[:i], ab.Books[i+1:]...)
				//ab.Books = ab.Books[:len(ab.Books)-1]
				// memory safety
				break
			}
		}
		fmt.Println(ab)
		dataHandler.AuthorList[dataHandler.SmStr(author.Name)] = ab
	}

	// setting new authors

	for _, author := range updBook.Authors {
		name := author.Name
		_, ok := dataHandler.AuthorList[dataHandler.SmStr(name)]
		fmt.Println("------------------")
		var ab dataHandler.AuthorBooks
		if ok == false {
			ab.Author = author
			ab.Books = append(ab.Books, updBook.ISBN)
			dataHandler.AuthorList[dataHandler.SmStr(name)] = ab
			continue
		}
		ab = dataHandler.AuthorList[dataHandler.SmStr(name)]
		ab.Books = append(ab.Books, updBook.ISBN)
		dataHandler.AuthorList[dataHandler.SmStr(name)] = ab
	}
	dataHandler.BookList[ISBN] = updBook
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Data added successfully"))
	if err != nil {
		http.Error(w, "Can not Write Data", http.StatusInternalServerError)
		return
	}
	fmt.Println("hello update func")

}

func GetSingleBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ISBN string
	ISBN = chi.URLParam(r, "ISBN")
	if len(ISBN) == 0 {
		http.Error(w, "ISBN is wrong", http.StatusBadRequest)
		return
	}
	_, ok := dataHandler.BookList[ISBN]
	if ok == false {
		http.Error(w, "ISBN is wrong", http.StatusBadRequest)
		return
	}
	var sb dataHandler.Book
	sb = dataHandler.BookList[ISBN]
	err := json.NewEncoder(w).Encode(sb)
	fmt.Println(sb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetAuthors(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(dataHandler.AuthorList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetSingleAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var AuthorName string
	AuthorName = chi.URLParam(r, "AuthorName")
	AuthorName = dataHandler.SmStr(AuthorName)
	fmt.Println(AuthorName)
	if len(AuthorName) == 0 {
		http.Error(w, "AuthorName is wrong", http.StatusBadRequest)
		return
	}
	_, ok := dataHandler.AuthorList[AuthorName]
	if ok == false {
		http.Error(w, "AuthorName is wrong", http.StatusBadRequest)
		return
	}
	var sa dataHandler.AuthorBooks
	sa = dataHandler.AuthorList[dataHandler.SmStr(AuthorName)]
	err := json.NewEncoder(w).Encode(sa)
	fmt.Println(sa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func RunServer(Port int) {
	dataHandler.Init()
	authHandler.InitToken()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)
	r.Group(func(r chi.Router) {
		r.Route("/books", func(r chi.Router) {
			r.Get("/", GetBooks)
			r.Get("/general", BookGeneralized)
			r.Get("/get/{ISBN}", GetSingleBook)
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(authHandler.TokenAuth))
				r.Use(jwtauth.Authenticator)
				r.Post("/", NewBook)
				r.Delete("/{ISBN}", DeleteBook)
				r.Put("/{ISBN}", UpdateBook)
			})
		})
		r.Route("/authors", func(r chi.Router) {
			r.Get("/", GetAuthors)
			r.Get("/{AuthorName}", GetSingleAuthor)
		})
	})

	if err := http.ListenAndServe(":"+strconv.Itoa(Port), r); err != nil {
		log.Fatalln(err)
	}
}
