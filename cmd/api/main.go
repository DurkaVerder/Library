package main

import (
	"library/internal/db"
	"library/internal/server"
)

func main() {
	db.InitDataBase()
	defer db.DataBaseClose()

	server.StartServer()
}
