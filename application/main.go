package main

import (
	"library/utils/db"
	"library/utils/sv"
)

func main() {
	db.InitDataBase()
	defer db.DB.Close()

	sv.StartServer()
}
