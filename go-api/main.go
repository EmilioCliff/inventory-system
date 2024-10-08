package main

import (
	"context"
	"fmt"
	"os"

	"github.com/EmilioCliff/inventory-system/api"
	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/EmilioCliff/inventory-system/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := utils.ReadConfig(".")
	if err != nil {
		log.Fatal().Msgf("Could not log config file: %s", err)
	}

	// runMigration("file://db/migration", config.DB_SOURCE_DEVELOPMENT)
	conn, err := pgxpool.New(context.Background(), config.DB_SOURCE_DEVELOPMENT)
	if err != nil {
		log.Fatal().Msgf("Couldnt connect to db: %s", err)
	}

	emailSender := utils.NewGmailSender(config.EMAIL_SENDER_NAME, config.EMAIL_SENDER_ADDRESS, config.EMAIL_SENDER_PASSWORD)

	store := db.NewStore(conn)
	redisOpt := asynq.RedisClientOpt{
		// Addr: "redis-inventory:6379",
		Addr:     config.REDIS_URI,
		Password: config.REDIS_PASSWORD,
		DB:       1,
	}

	redisClient := redis.NewClient(&redis.Options{
		// Addr: "redis-inventory:6379",
		Addr:     config.REDIS_URI,
		Password: config.REDIS_PASSWORD,
		DB:       2,
	})

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	server, err := api.NewServer(config, store, *emailSender, taskDistributor, redisClient)
	if err != nil {
		log.Fatal().Msgf("Couldnt create new server: %s", err)
	}

	accessToken, err := server.GeneratePythonToken("pythonApp")
	fmt.Println(accessToken)

	port := os.Getenv("PORT")

	go runRedisTaskProcessor(redisOpt, *store, *emailSender, config, taskDistributor)
	log.Info().Msgf("starting server at port: %s", port)
	err = server.Start(fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatal().Msgf("Couldnot start server: %s", err)
	}
}

func runRedisTaskProcessor(
	redisOpt asynq.RedisClientOpt,
	store db.Store,
	sender utils.GmailSender,
	config utils.Config,
	distributor worker.TaskDistributor,
) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, sender, config, distributor)
	log.Info().Msg("Start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Msgf("could not start task processor: %s", err)
	}

	// here
	// ctx := context.Background()
	// opts := []asynq.Option{
	// 	asynq.MaxRetry(2),
	// 	asynq.ProcessIn(10 * time.Second),
	// 	asynq.Queue(worker.QueueLow),
	// }
	// if err := distributor.DistributeTakeAndSendDBsnapshots(ctx, "word", opts...); err != nil {
	// 	log.Fatal().Msgf("Failed to distribute task: %s", err)
	// }
}

func runMigration(mirationUrl string, db_source string) {
	migration, err := migrate.New(mirationUrl, db_source)
	if err != nil {
		log.Fatal().Msgf("Failed to load migration: %s", err)
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("Failed to run migrate up: %s", err)
	}

	log.Info().Msg("Migration Successfull")
}
