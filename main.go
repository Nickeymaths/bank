package main

import (
	"database/sql"
	"log"

	"github.com/Nickeymaths/bank/api"
	db "github.com/Nickeymaths/bank/db/sqlc"
	"github.com/Nickeymaths/bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load server config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Fail to start server: ", err)
	}
}
