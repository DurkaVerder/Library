package server

import (
	"library/internal/authentication"
	"library/internal/handler"
	"log"
	"net/http"
)

func StartServer() {
	port := ":8080"
	http.HandleFunc("/api/v1/books", handler.HandleRequest)
	http.HandleFunc("/api/v1/books/search", handler.HandlePaginationSort)
	http.HandleFunc("/api/v1/book/", handler.HandleRequestWithId)
	http.HandleFunc("/api/v1/books/sort", handler.HandleRequestSort)
	http.HandleFunc("/api/v1/login", authentication.HandleLogin)
	http.HandleFunc("/api/v1/register", authentication.HandleRegister)

	log.Println("Starting server")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}

}
