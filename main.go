package main

import (
	"BOOKS-LIST/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var books []models.Book
var db *sql.DB

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	// handleRequests()
	// router := mux.NewRouter()
	/*
		db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/library")

		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query("INSERT INTO books VALUES (2,'Golang GoRoutines','Mr.GoRoutines','2010') ")
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	*/

	books = append(books,
		Book{ID: 1, Title: "Golang Pointers", Author: "Mr. Golang", Year: "2010"},
		Book{ID: 2, Title: "Golang GoRoutines", Author: "Mr. GoRoutines", Year: "2011"},
		Book{ID: 3, Title: "Golang Routers", Author: "Mr. Routers", Year: "2012"},
		Book{ID: 4, Title: "Golang Concurrency", Author: "Mr. Concurrency", Year: "2013"},
	)

	http.HandleFunc("/articles", getBooks)

	/*
		router.HandleFunc("/books", getBooks).Methods("GET")
		router.HandleFunc("/books/{id}", getBook).Methods("GET")
		router.HandleFunc("/books", addBook).Methods("POST")
		router.HandleFunc("/books", updateBook).Methods("PUT")
		router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")
	*/
	log.Fatal(http.ListenAndServe(":2222", nil))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	// rows, err := db.Query("SELECT * FROM books")
	// log.Println(err)
	// defer rows.Close()

	// for rows.Next() {
	// 	var book Book
	// 	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	books = append(books, book)
	// }
	// log.Println("helllo")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	rows := db.QueryRow("SELECT * FROM books where id=$1", id)
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	log.Println(err)
	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)
	json.NewEncoder(w).Encode(books)

	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}

	json.NewEncoder(w).Encode(books)

}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, item := range books {
		if item.ID == id {
			books = append(books[:i], books[i+1:]...)
		}
	}

	json.NewEncoder(w).Encode(books)
}
