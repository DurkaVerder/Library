package sv

import (
	"library/authentication"
	"library/utils/db"
	"log"
	"net/http"
	"strconv"
)

func StartServer() {
	port := ":8080"
	http.HandleFunc("/books", handleRequest)
	http.HandleFunc("/book/", handleRequestWithId)
	http.HandleFunc("/books/sort", handleRequestSort)
	http.HandleFunc("/login", authentication.HandleLogin)
	log.Println("Starting server")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		db.GetBooks(w)
	case http.MethodPost:
		db.PostBook(w, r)
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
		db.GetIdBook(w, id)
	case http.MethodDelete:
		db.DeleteBook(w, id)
	case http.MethodPut:
		db.PutBook(w, r, id)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func handleRequestSort(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		db.SortBooks(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
