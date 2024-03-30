package main

import (
	"context"
	"fmt"

	"github.com/EmilioCliff/inventory-system/api"
	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/EmilioCliff/inventory-system/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := utils.ReadConfig(".")
	if err != nil {
		log.Fatal().Msgf("Could not log config file: %s", err)
	}
	conn, err := pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatal().Msgf("Couldnt connect to db: %s", err)
	}

	emailSender := utils.NewGmailSender(config.EMAIL_SENDER_NAME, config.EMAIL_SENDER_ADDRESS, config.EMAIL_SENDER_PASSWORD)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr:     config.REDIS_ADDRESS,
		Password: config.REDIS_PASSWORD,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	server, err := api.NewServer(config, store, *emailSender, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("Couldnt create new server: %s", err)
	}

	accessToken, err := server.GeneratePythonToken("pythonApp")
	fmt.Println(accessToken)

	go runRedisTaskProcessor(redisOpt, *store, *emailSender, config)
	err = server.Start(config.SERVER_ADDRESS)
	log.Info().Msgf("starting server at port: %s", config.SERVER_ADDRESS)
	if err != nil {
		log.Fatal().Msgf("Couldnot start server: %s", err)
	}
}

func runRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, sender utils.GmailSender, config utils.Config) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, sender, config)
	log.Info().Msg("Start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Msgf("could not start task processor: %s", err)
	}
}
