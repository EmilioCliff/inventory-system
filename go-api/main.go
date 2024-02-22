package main

import (
	"context"
	"fmt"
	"log"

	"github.com/EmilioCliff/inventory-system/api"
	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := utils.ReadConfig(".")
	if err != nil {
		log.Fatal("Could not log config file: ", err)
	}
	conn, err := pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatal("Couldnt connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldnt create new server: ", err)
	}

	accessToken, err := server.GeneratePythonToken("pythonApp")
	fmt.Println(accessToken)

	err = server.Start(config.SERVER_ADDRESS)
	if err != nil {
		log.Fatal("Couldnot start server: ", err)
	}
}
