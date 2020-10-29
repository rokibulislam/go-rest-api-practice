package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/books-list/controllers"
	"github.com/books-list/driver"
	"github.com/books-list/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}

var Articles []Article
var books []models.Book

func handleRequests() {
	db := driver.ConnectDB()
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	controller := controllers.Controller{}
	myRouter.HandleFunc("/books", controller.GetBooks(db)).Methods("GET")
	myRouter.HandleFunc("/books/{id}", controller.GetBook(db)).Methods("GET")
	myRouter.HandleFunc("/books", controller.AddBook(db)).Methods("POST")
	myRouter.HandleFunc("/books", controller.UpdateBook(db)).Methods("PUT")
	myRouter.HandleFunc("/books/{id}", controller.RemoveBook(db)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":30000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")

	Articles = []Article{
		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}

	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/library")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Books")
	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/library")
	defer db.Close()

	rows, err := db.Query("SELECT * FROM books")
	log.Println(err)
	defer rows.Close()

	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			panic(err.Error())
		}
		books = append(books, book)
		log.Println(book.Title)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Book")
	var book models.Book
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/library")
	defer db.Close()
	rows := db.QueryRow("SELECT * FROM books where id= ?", id)
	b := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	log.Println(err)
	log.Println(b)
	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit:Add Book")
	var book models.Book
	var bookID int
	json.NewDecoder(r.Body).Decode(&book)
	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/library")
	defer db.Close()
	err = db.QueryRow("INSERT INTO books (title,author,year) VALUES ($1,$2,$3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)
	log.Println(err)
	books = append(books, book)
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
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

func removeBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit:DELETE Book")

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	result, _ := db.Exec("DELETE books where id=$1 RETURNING id;", id)
	rowdeleted, _ := result.RowsAffected()
	json.NewEncoder(w).Encode(rowdeleted)
}
