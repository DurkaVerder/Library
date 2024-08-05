package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDataBase() {
	openDB := "user=postgres password=durka dbname=storeBooks sslmode=disable host=host.docker.internal port=5432"
	var err error
	db, err = sql.Open("postgres", openDB)
	if err != nil {
		log.Fatal("Error open database: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting database: ", err)
	}
}

func DataBaseClose() {
	db.Close()
}

func GetBD() *sql.DB {
	return db
}
