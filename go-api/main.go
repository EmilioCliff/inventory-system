package main

import (
	"context"
	"fmt"
	"os"

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
	conn, err := pgxpool.New(context.Background(), os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal().Msgf("Couldnt connect to db: %s", err)
	}

	emailSender := utils.NewGmailSender(os.Getenv("EMAIL_SENDER_NAME"), os.Getenv("EMAIL_SENDER_ADDRESS"), os.Getenv("EMAIL_SENDER_PASSWORD"))

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	server, err := api.NewServer(config, store, *emailSender, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("Couldnt create new server: %s", err)
	}

	accessToken, err := server.GeneratePythonToken("pythonApp")
	fmt.Println(accessToken)

	go runRedisTaskProcessor(redisOpt, *store, *emailSender, config, taskDistributor)
	err = server.Start(os.Getenv("SERVER_ADDRESS"))
	log.Info().Msgf("starting server at port: %s", os.Getenv("SERVER_ADDRESS"))
	if err != nil {
		log.Fatal().Msgf("Couldnot start server: %s", err)
	}
}

func runRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, sender utils.GmailSender, config utils.Config, distributor worker.TaskDistributor) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, sender, config, distributor)
	log.Info().Msg("Start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Msgf("could not start task processor: %s", err)
	}
}

// func runConfig() (utils.Config, error) {
// 	config, err := utils.ReadConfig(".")
// 	if err != nil {
// 		return config, err
// 	}
// 	// Other initialization logic goes here
// 	return config, nil
// }
