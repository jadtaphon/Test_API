package main

import (
	"fmt"
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	//"golang.org/x/net/context"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

var server = "localhost"
var port = 1433
var user = "sa"
var password = "adasoft"
var database = "TestJa"

var db *sql.DB
var err error

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)

	db, err = sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/allbook", getAllBooks).Methods("GET")
	router.HandleFunc("/searchbook/{id}", searchBook).Methods("GET")
	router.HandleFunc("/createbook", createBook).Methods("POST")
	router.HandleFunc("/updatebook/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/deletebook/{id}", deleteBook).Methods("DELETE")

	log.Printf("Running at localhost:8080\n")
	http.ListenAndServe(":8080", router)
}

type book struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Year string `json:"year"`
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM book")
	if err != nil {
		log.Fatal("Database SELECT failed:", err.Error())
	}
	defer rows.Close()

	var books []book
	for rows.Next() {
		var b book
		err := rows.Scan(&b.Id, &b.Name, &b.Year)
		if err != nil {
			log.Fatal("Scan failed:", err.Error())
		}
		books = append(books, b)
	}

	json.NewEncoder(w).Encode(books)
}

func searchBook(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	rows, err := db.Query("SELECT * FROM book WHERE id=?", param["id"])
	if err != nil {
		log.Fatal("Database SELECT failed:", err.Error())
	}
	defer rows.Close()

	var b book
	for rows.Next() {
		err := rows.Scan(&b.Id, &b.Name, &b.Year)
		if err != nil {
			log.Fatal("Scan failed:", err.Error())
		}
	}

	json.NewEncoder(w).Encode(b)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var b book
	json.NewDecoder(r.Body).Decode(&b)

	stmt, err := db.Prepare("INSERT INTO book(name, year) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}

	_, err = stmt.Exec(&b.Id, &b.Name, &b.Year)
	if err != nil {
		log.Fatal("Database INSERT failed:", err.Error())
	}
	defer stmt.Close()

	w.WriteHeader(http.StatusCreated)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	var b book
	json.NewDecoder(r.Body).Decode(&b)

	stmt, err := db.Prepare("UPDATE book SET name=?, year=? WHERE id=?")
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}

	_, err = stmt.Exec(&b.Id, &b.Name, &b.Year, param["id"])
	if err != nil {
		log.Fatal("Database UPDATE failed:", err.Error())
	}
	defer stmt.Close()

	w.WriteHeader(http.StatusOK)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM book WHERE id=?")
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}

	_, err = stmt.Exec(param["id"])
	if err != nil {
		log.Fatal("Database DELETE failed:", err.Error())
	}
	defer stmt.Close()

	w.WriteHeader(http.StatusOK)
}