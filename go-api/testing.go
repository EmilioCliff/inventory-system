package main

import (
	"context"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/EmilioCliff/inventory-system/reports"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func maini() {
	config, err := utils.ReadConfig(".")
	if err != nil {
		log.Fatal().Msgf("Could not log config file: %s", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DB_SOURCE_DEVELOPMENT)
	if err != nil {
		log.Fatal().Msgf("Couldnt connect to db: %s", err)
	}

	store := db.NewStore(conn)

	fromDate, err := time.Parse("2006-01-02", "2024-08-01")
	if err != nil {
		log.Fatal().Msgf("Couldnt parse time: %s", err)
	}

	toDate, err := time.Parse("2006-01-02", "2024-08-31")
	if err != nil {
		log.Fatal().Msgf("Couldnt parse time: %s", err)
	}
	reportMaker := reports.NewReportMaker(store)

	_, err = reportMaker.GenerateUserExcel(context.Background(), reports.ReportsPayload{
		FromDate: fromDate,
		ToDate:   toDate,
	})
	if err != nil {
		log.Fatal().Msgf("This is some error: %s", err)
	} else {
		log.Printf("successful excel creation")
	}

	_, err = reportMaker.GenerateManagerReports(context.Background(), reports.ReportsPayload{
		FromDate: fromDate,
		ToDate:   toDate,
	})
	if err != nil {
		log.Fatal().Msgf("This is some error from manager: %s", err)
	} else {
		log.Printf("successful excel creation for manager")
	}
}
