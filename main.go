package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"strings"
	"time"
)

var jwtkey = []byte("this_is_a_secret_key")
var tokenAuth *jwtauth.JWTAuth
var tokenString string
var token jwt.Token

type Author struct {
	Name      string `json:"name,omitempty"`
	BookCount string `json:"book_count"`
	Age       string `json:"age"`
}

type Book struct {
	Name    string   `json:"book_name,omitempty"`
	Authors []Author `json:"authors"`
	ISBN    string   `json:"isbn,omitempty"`
	Genre   string   `json:"genre"`
	Pub     string   `json:"publisher"`
}

type AuthorBooks struct {
	Author `json:"author"`
	Books  []string `json:"books"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BookDB map[string]Book
type AuthorDB map[string]AuthorBooks
type CredDB map[string]string

var BookList BookDB
var AuthorList AuthorDB
var CredList CredDB

func CapToSmall(s string) string {
	return strings.ToLower(s)
}

func RmSpace(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func SmStr(s string) string {
	return CapToSmall(RmSpace(s))
}

func Init() {
	author1 := Author{
		Name:      "temp author 1",
		BookCount: "5",
		Age:       "45",
	}
	author2 := Author{
		Name:      "temp author 2",
		BookCount: "5",
		Age:       "45",
	}

	data1 := Book{
		Name: "temp book 1",
		Authors: []Author{
			author1,
			author2,
		},
		ISBN:  "ISBN1",
		Genre: "Fiction",
		Pub:   "Demo",
	}

	data2 := Book{
		Name: "temp book 2",
		Authors: []Author{
			author1,
		},
		ISBN:  "ISBN2",
		Genre: "Fiction",
		Pub:   "Demo",
	}

	data3 := Book{
		Name: "temp book 3",
		Authors: []Author{
			author2,
		},
		ISBN:  "ISBN3",
		Genre: "Fiction",
		Pub:   "Demo",
	}

	User := Credentials{
		Username: "user",
		Password: "pass",
	}

	BookList = make(BookDB)
	AuthorList = make(AuthorDB)
	CredList = make(CredDB)

	var ab1 AuthorBooks
	ab1.Author = author1
	ab1.Books = append(ab1.Books, data1.ISBN)
	ab1.Books = append(ab1.Books, data2.ISBN)

	var ab2 AuthorBooks
	ab2.Author = author2
	ab2.Books = append(ab2.Books, data1.ISBN)
	ab2.Books = append(ab2.Books, data3.ISBN)

	AuthorList[SmStr(author1.Name)] = ab1
	AuthorList[SmStr(author2.Name)] = ab2

	BookList[data1.ISBN] = data1
	BookList[data2.ISBN] = data2
	BookList[data3.ISBN] = data3

	CredList[User.Username] = User.Password
	tokenAuth = jwtauth.New(string(jwa.HS256), jwtkey, nil)
	return
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(BookList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func BookGeneralized(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	var readString []string
	for _, str := range BookList {
		readString = append(readString, str.Name)
	}
	resp := strings.Join(readString, "\n")
	if resp == "" {
		http.Error(w, "No Books found", http.StatusNotFound)
		return
	}
	w.Write([]byte(resp))
}

func NewBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Cannot Decode the data", http.StatusBadRequest)
		return
	}
	if len(book.Name) == 0 || len(book.ISBN) == 0 || len(book.Authors) == 0 {
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

	BookList[book.ISBN] = book

	for _, author := range book.Authors {
		name := author.Name
		_, ok := AuthorList[SmStr(name)]
		var ab AuthorBooks
		if ok == false {
			ab.Author = author
			ab.Books = append(ab.Books, book.ISBN)
			AuthorList[SmStr(name)] = ab
			continue
		}
		ab = AuthorList[SmStr(name)]
		ab.Books = append(ab.Books, book.ISBN)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data added successfully"))
	BookList[book.ISBN] = book
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
	_, ok := BookList[ISBN]
	if ok == false {
		http.Error(w, "ISBN is wrong", http.StatusBadRequest)
		return
	}

	delete(BookList, ISBN)

	for _, data := range AuthorList {
		var ab AuthorBooks
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
	w.Write([]byte("Data Deleted successfully"))
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
	_, ok := BookList[ISBN]
	if ok == false {
		http.Error(w, "ISBN is wrong", http.StatusBadRequest)
		return
	}

	curBook := BookList[ISBN]
	var updBook Book

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
		var ab AuthorBooks
		ab = AuthorList[SmStr(author.Name)]
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
		AuthorList[SmStr(author.Name)] = ab
	}

	// setting new authors

	for _, author := range updBook.Authors {
		name := author.Name
		_, ok := AuthorList[SmStr(name)]
		fmt.Println("------------------")
		var ab AuthorBooks
		if ok == false {
			ab.Author = author
			ab.Books = append(ab.Books, updBook.ISBN)
			AuthorList[SmStr(name)] = ab
			continue
		}
		ab = AuthorList[SmStr(name)]
		ab.Books = append(ab.Books, updBook.ISBN)
		AuthorList[SmStr(name)] = ab
	}
	BookList[ISBN] = updBook
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data Updated successfully"))
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
	_, ok := BookList[ISBN]
	if ok == false {
		http.Error(w, "ISBN is wrong", http.StatusBadRequest)
		return
	}
	var sb Book
	sb = BookList[ISBN]
	err := json.NewEncoder(w).Encode(sb)
	fmt.Println(sb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(AuthorList)

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
	AuthorName = SmStr(AuthorName)
	fmt.Println(AuthorName)
	if len(AuthorName) == 0 {
		http.Error(w, "AuthorName is wrong", http.StatusBadRequest)
		return
	}
	_, ok := AuthorList[AuthorName]
	if ok == false {
		http.Error(w, "AuthorName is wrong", http.StatusBadRequest)
		return
	}
	var sa AuthorBooks
	sa = AuthorList[SmStr(AuthorName)]
	err := json.NewEncoder(w).Encode(sa)
	fmt.Println(sa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	fmt.Println(cred)
	if err != nil {
		http.Error(w, "Can not Decode the data", http.StatusBadRequest)
		return
	}

	password, okay := CredList[cred.Username]

	if okay == false {
		http.Error(w, "Username do not exist", http.StatusBadRequest)
		return
	}

	if password != cred.Password {
		http.Error(w, "Password not matching", http.StatusBadRequest)
		return
	}
	fmt.Println("chekc")
	et := time.Now().Add(15 * time.Minute)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"aud": "Saurov Biswas",
		"exp": et.Unix(),
	})
	fmt.Println(tokenString)
	if err != nil {
		http.Error(w, "Can not generate jwt", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   tokenString,
		Expires: et,
	})
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
}

func main() {
	Init()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Post("/login", Login)
	r.Post("/logout", Logout)
	r.Group(func(r chi.Router) {
		r.Route("/books", func(r chi.Router) {
			r.Get("/", GetBooks)
			r.Get("/general", BookGeneralized)
			r.Get("/getbook/{ISBN}", GetSingleBook)
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
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

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}
