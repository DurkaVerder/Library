package main

import (
	"library/config/cfg"
	"library/internal/db"
	"library/internal/rd"
	"library/internal/server"
	"log"
)

func main() {
	err := cfg.InitCfg()
	if err != nil {
		log.Println("Error init config: ", err)
		return
	}

	db.InitDataBase()
	defer db.DataBaseClose()

	rd.InitRedis()
	defer rd.CloseRdb()

	server.StartServer()
}
