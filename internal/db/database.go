package db

import (
	"database/sql"
	"fmt"
	"library/config/cfg"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDataBase() {
	conf := cfg.Cfg
	openDB := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable host=%v port=%v", conf.DataBase.User, conf.DataBase.Password, conf.DataBase.Dbname, conf.DataBase.Host, conf.DataBase.Port)
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
