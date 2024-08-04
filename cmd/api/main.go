package main

import (
	"library/internal/db"
	"library/internal/rd"
	"library/internal/server"
)

func main() {
	db.InitDataBase()
	defer db.DataBaseClose()

	rd.InitRedis()
	defer rd.CloseRdb()

	server.StartServer()
}
