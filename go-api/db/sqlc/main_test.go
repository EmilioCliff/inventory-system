package db

import (
	"context"
	"log"
	"testing"

	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore *Queries
var testConnPool *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := utils.ReadConfig("../..")
	if err != nil {
		log.Fatal("Could not log config")
	}
	testConnPool, err = pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatal("Couldnot connect to db: ", err)
	}

	testStore = New(testConnPool)

	// os.Exit(m.Run())
}
