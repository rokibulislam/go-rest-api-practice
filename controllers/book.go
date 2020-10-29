package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/books-list/models"
	"github.com/books-list/repository/bookRepository"
	"github.com/gorilla/mux"
)

var books []models.Book

type Controller struct{}

func (c Controller) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint Hit: Books")
		var book models.Book
		// var error models.Error

		books = []models.Book{}
		bookRepo := bookRepository.BookRepository{}
		books, _ := bookRepo.GetBooks(db, book, books)

		// rows, err := db.Query("SELECT * FROM books")
		// log.Println(err)
		// defer rows.Close()
		// for rows.Next() {
		// 	var book models.Book
		// 	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		// 	if err != nil {
		// 		panic(err.Error())
		// 	}
		// 	books = append(books, book)
		// 	log.Println(book.Title)
		// }

		json.NewEncoder(w).Encode(books)
	}
}

func (c Controller) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint Hit: Book")
		var book models.Book
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])
		rows := db.QueryRow("SELECT * FROM books where id= ?", id)
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		log.Println(err)
		json.NewEncoder(w).Encode(book)
	}
}

func (c Controller) AddBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint Hit:Add Book")
		var book models.Book
		var bookID int
		json.NewDecoder(r.Body).Decode(&book)
		err := db.QueryRow("INSERT INTO books (title,author,year) VALUES ($1,$2,$3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)
		log.Println(err)
		books = append(books, book)
		json.NewEncoder(w).Encode(books)
	}
}

func (c Controller) UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint Hit:Update Book")
		var book models.Book
		json.NewDecoder(r.Body).Decode(&book)
		db, _ := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/library")
		defer db.Close()
		result, _ := db.Exec("UPDATE  books set title=$1 author=$2 year=$3 where id=$4 RETURNING id;", book.Title, book.Author, book.Year, book.ID)
		// log.Println(er)
		rowupdated, _ := result.RowsAffected()
		json.NewEncoder(w).Encode(rowupdated)
	}
}

func (c Controller) RemoveBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Endpoint Hit:DELETE Book")
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])
		result, _ := db.Exec("DELETE books where id=$1 RETURNING id;", id)
		rowdeleted, _ := result.RowsAffected()
		json.NewEncoder(w).Encode(rowdeleted)
	}
}
