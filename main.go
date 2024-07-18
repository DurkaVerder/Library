package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Book struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Details string `json:"details"`
	Author  string `json:"author"`
}

type TypeSort struct {
	Name string `json:"name"`
}

var db *sql.DB

func main() {
	initDataBase()
	defer db.Close()

	startServer()

}

func initDataBase() {
	openDB := "user=postgres password=durka dbname=storeBooks sslmode=disable"
	var err error
	db, err = sql.Open("postgres", openDB)
	if err != nil {
		log.Fatal("Error open database: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting database: ", err)
	}

}

func startServer() {
	port := ":8080"
	http.HandleFunc("/books", handleRequest)
	http.HandleFunc("/book/", handleRequestWithId)
	http.HandleFunc("/books/sort", handleRequestSort)
	log.Println("Starting server")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBooks(w)
	case http.MethodPost:
		postBook(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func handleRequestWithId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/book/"):])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getIdBook(w, id)
	case http.MethodDelete:
		deleteBook(w, id)
	case http.MethodPut:
		putBook(w, r, id)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func handleRequestSort(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		sortBooks(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func getBooks(w http.ResponseWriter) {
	getSQL := `SELECT * FROM books`
	listBooks := make([]Book, 0)
	rows, err := db.Query(getSQL)
	if err != nil {
		log.Println("Error select-request")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.Id, &b.Name, &b.Details, &b.Author); err != nil {
			log.Println("Error scan rows: ", err)
			return
		}
		listBooks = append(listBooks, b)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows: ", err)
		return
	}

	json.NewEncoder(w).Encode(listBooks)
	w.WriteHeader(http.StatusOK)

}

func postBook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error read request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var b Book
	if err := json.Unmarshal(body, &b); err != nil {
		log.Println("Error unmarshal json: ", err)
	}

	postSQL := `INSERT INTO books (name, details, author) VALUES ($1, $2, $3)`
	if _, err := db.Exec(postSQL, b.Name, b.Details, b.Author); err != nil {
		log.Println("Error post-request: ", err)
		return
	}

	json.NewEncoder(w).Encode(b)
	w.WriteHeader(http.StatusCreated)

}

func deleteBook(w http.ResponseWriter, id int) {
	deleteSQL := `DELETE FROM books WHERE id = $1`
	if _, err := db.Exec(deleteSQL, id); err != nil {
		log.Println("Error delete-request: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

func putBook(w http.ResponseWriter, r *http.Request, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error read request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var b Book
	if err := json.Unmarshal(body, &b); err != nil {
		log.Println("Error unmarshal json: ", err)
	}

	putSQL := `UPDATE books SET name = $1, details = $2, author = $3 WHERE id = $4`
	if _, err := db.Exec(putSQL, b.Name, b.Details, b.Author, id); err != nil {
		log.Println("Error put-request: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func getIdBook(w http.ResponseWriter, id int) {
	getSQL := `SELECT * FROM books WHERE id = $1`

	row := db.QueryRow(getSQL, id)

	var b Book
	if err := row.Scan(&b.Id, &b.Name, &b.Details, &b.Author); err != nil {
		log.Println("Error scan rows: ", err)
		return
	}

	if err := row.Err(); err != nil {
		log.Println("Error in rows: ", err)
		return
	}

	json.NewEncoder(w).Encode(b)
	w.WriteHeader(http.StatusOK)
}

func sortBooks(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error read request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var s TypeSort
	if err := json.Unmarshal(body, &s); err != nil {
		log.Println("Error unmarshal json: ", err)
	}

	getSQL := fmt.Sprintf(`SELECT * FROM books ORDER BY %s`, s.Name)
	listBooks := make([]Book, 0)
	rows, err := db.Query(getSQL)
	if err != nil {
		log.Println("Error select-request")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.Id, &b.Name, &b.Details, &b.Author); err != nil {
			log.Println("Error scan rows: ", err)
			return
		}
		listBooks = append(listBooks, b)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows: ", err)
		return
	}

	json.NewEncoder(w).Encode(listBooks)
	w.WriteHeader(http.StatusOK)
}
