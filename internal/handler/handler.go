package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"library/internal/db"
	"library/internal/jwt"
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

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if err := jwt.CheckJWT(w, r); err == 1 {
		log.Println("Error checkJWT")
		return
	}

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

func HandleRequestWithId(w http.ResponseWriter, r *http.Request) {
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

func HandleRequestSort(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		sortBooks(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func HandlePaginationSort(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("limit")
	genreParam := r.URL.Query().Get("genre")

	limit := 10

	if limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error to convert limit param: ", err)
			return
		}
	}
	var rows *sql.Rows
	var err error
	selectSQL := `SELECT * FROM books WHERE details = $1 LIMIT $2`
	if genreParam == "" {
		selectSQL = `SELECT * FROM books LIMIT $1`
		rows, err = db.GetBD().Query(selectSQL, limit)
	} else {
		rows, err = db.GetBD().Query(selectSQL, genreParam, limit)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error select-request: ", err)
		return
	}
	defer rows.Close()

	books := make([]Book, 0, limit)
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.Id, &b.Name, &b.Details, &b.Author); err != nil {
			log.Println("Error scan rows: ", err)
			return
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows: ", err)
		return
	}

	json.NewEncoder(w).Encode(books)
	w.WriteHeader(http.StatusOK)

}

func getBooks(w http.ResponseWriter) {
	getSQL := `SELECT * FROM books`
	listBooks := make([]Book, 0)
	rows, err := db.GetBD().Query(getSQL)
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
	if _, err := db.GetBD().Exec(postSQL, b.Name, b.Details, b.Author); err != nil {
		log.Println("Error post-request: ", err)
		return
	}

	json.NewEncoder(w).Encode(b)
	w.WriteHeader(http.StatusCreated)

}

func deleteBook(w http.ResponseWriter, id int) {
	deleteSQL := `DELETE FROM books WHERE id = $1`
	if _, err := db.GetBD().Exec(deleteSQL, id); err != nil {
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
	if _, err := db.GetBD().Exec(putSQL, b.Name, b.Details, b.Author, id); err != nil {
		log.Println("Error put-request: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func getIdBook(w http.ResponseWriter, id int) {
	getSQL := `SELECT * FROM books WHERE id = $1`

	row := db.GetBD().QueryRow(getSQL, id)

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
	rows, err := db.GetBD().Query(getSQL)
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
