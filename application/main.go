package main

import (
	"library/utils/db"
	"library/utils/sv"
)

// @title Bookstore API
// @version 1.0
// @description This is a sample server for a bookstore.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.email vrrrr227@gmail.com

// @license.name MIT
// @license.url http://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

func main() {
	db.InitDataBase()
	defer db.DB.Close()

	sv.StartServer()
}
