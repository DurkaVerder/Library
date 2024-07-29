package sv

import (
	"library/authentication"
	_ "library/docs"
	"library/utils/db"
	"library/utils/jwt"
	"log"
	"net/http"
	"strconv"

	httpSwagger "github.com/swaggo/http-swagger"
)

func StartServer() {
	port := ":8080"
	http.HandleFunc("/api/v1/books", handleRequest)
	http.HandleFunc("/api/v1/books/search", db.HandlePaginationSort)
	http.HandleFunc("/api/v1/book/", handleRequestWithId)
	http.HandleFunc("/api/v1/books/sort", handleRequestSort)
	http.HandleFunc("/api/v1/login", authentication.HandleLogin)
	http.HandleFunc("/api/v1/register", authentication.HandleRegister)
	http.HandleFunc("/api/v1/swagger/", httpSwagger.WrapHandler)
	log.Println("Starting server")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}

}

// HandleRequest godoc
// @Summary Check JWT and choose method for request
// @Description This endpoint checks the JWT token and determines the appropriate method for handling the request.
// @Tags Handle
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 400 {object} string "Bad Request"

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if err := jwt.CheckJWT(w, r); err == 1 {
		log.Println("Error checkJWT")
		return
	}

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
