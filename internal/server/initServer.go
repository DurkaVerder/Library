package server

import (
	"library/config/cfg"
	"library/internal/authentication"
	"library/internal/handler"
	"log"
	"net/http"
)

func StartServer() {
	conf := cfg.Cfg
	http.HandleFunc("/api/v1/books", handler.HandleRequest)
	http.HandleFunc("/api/v1/books/search", handler.HandlePaginationSort)
	http.HandleFunc("/api/v1/book/", handler.HandleRequestWithId)
	http.HandleFunc("/api/v1/books/sort", handler.HandleRequestSort)
	http.HandleFunc("/api/v1/login", authentication.HandleLogin)
	http.HandleFunc("/api/v1/register", authentication.HandleRegister)

	log.Println("Starting server")
	if err := http.ListenAndServe(conf.Server.Port, nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}

}
