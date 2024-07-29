package db

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

type User struct {
	User_id  int    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Rule     string `json:"rule"`
}

var DB *sql.DB

func InitDataBase() {
	openDB := "user=postgres password=durka dbname=storeBooks sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", openDB)
	if err != nil {
		log.Fatal("Error open database: ", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting database: ", err)
	}
}

func GetBooks(w http.ResponseWriter) {
	getSQL := `SELECT * FROM books`
	listBooks := make([]Book, 0)
	rows, err := DB.Query(getSQL)
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

func PostBook(w http.ResponseWriter, r *http.Request) {
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
	if _, err := DB.Exec(postSQL, b.Name, b.Details, b.Author); err != nil {
		log.Println("Error post-request: ", err)
		return
	}

	json.NewEncoder(w).Encode(b)
	w.WriteHeader(http.StatusCreated)

}

func DeleteBook(w http.ResponseWriter, id int) {
	deleteSQL := `DELETE FROM books WHERE id = $1`
	if _, err := DB.Exec(deleteSQL, id); err != nil {
		log.Println("Error delete-request: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

func PutBook(w http.ResponseWriter, r *http.Request, id int) {
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
	if _, err := DB.Exec(putSQL, b.Name, b.Details, b.Author, id); err != nil {
		log.Println("Error put-request: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func GetIdBook(w http.ResponseWriter, id int) {
	getSQL := `SELECT * FROM books WHERE id = $1`

	row := DB.QueryRow(getSQL, id)

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

func SortBooks(w http.ResponseWriter, r *http.Request) {

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
	rows, err := DB.Query(getSQL)
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

func CheckLogin(login string) bool {
	check := `SELECT login FROM users WHERE login = $1`
	row := DB.QueryRow(check, login)
	var existingLogin string
	if err := row.Scan(&existingLogin); err == sql.ErrNoRows {
		return true
	} else if err != nil {
		log.Println("Error checking login")
		return false
	}
	return false
}

func AddUser(login, password string) error {
	addUser := `INSERT INTO users (login, password, rule) VALUES ($1, $2, $3)`
	if _, err := DB.Exec(addUser, login, password, "default"); err != nil {
		log.Println("Error add user: ", err)
		return err
	}
	return nil

}

func CheckExistUser(login, password string) bool {
	check := `SELECT user_id, login, password, rule FROM users WHERE login = $1 AND password = $2`
	row := DB.QueryRow(check, login, password)

	var u User
	if err := row.Scan(&u.User_id, &u.Login, &u.Password, &u.Rule); err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Println("Error scan row: ", err)
		return false
	}
	return true
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
		rows, err = DB.Query(selectSQL, limit)
	} else {
		rows, err = DB.Query(selectSQL, genreParam, limit)
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
