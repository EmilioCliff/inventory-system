package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TakeAndSendDBsnpashotsTask = "task:take_and_send_dbsnapshots"

func (distributor *RedisTaskDistributor) DistributeTakeAndSendDBsnapshots(ctx context.Context, word string, opts ...asynq.Option) error {
	jsonTaskPayload, err := json.Marshal(word)
	if err != nil {
		return fmt.Errorf("failed to marshal take and send bdsnapshots payload: %w", err)
	}

	task := asynq.NewTask(TakeAndSendDBsnpashotsTask, jsonTaskPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task take and send bdsnapshots: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("body", "database snapshot").
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTakeAndSendDBsnapshots(ctx context.Context, task *asynq.Task) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	snapshotSchemaFilename := timestamp + "_schema_snapshot.sql"
	snapshotDataAndSchemaFilename := timestamp + "_data_schema_snapshot.sql"

	cmd1 := exec.Command("docker", "exec", "postgres3", "/bin/bash", "-c", fmt.Sprintf("pg_dump -U root -d inventorydb -s -f %s%s", processor.config.POSTGRES_SNAPSHOTS, snapshotSchemaFilename))
	cmd2 := exec.Command("docker", "exec", "postgres3", "/bin/bash", "-c", fmt.Sprintf("pg_dump -U root -d inventorydb -f %s%s --create", processor.config.POSTGRES_SNAPSHOTS, snapshotDataAndSchemaFilename))
	cmd3 := exec.Command("docker", "cp", fmt.Sprintf("postgres3:%s%s", processor.config.POSTGRES_SNAPSHOTS, snapshotSchemaFilename), fmt.Sprintf("%s", processor.config.HOST_SNAPSHOTS))
	cmd4 := exec.Command("docker", "cp", fmt.Sprintf("postgres3:%s%s", processor.config.POSTGRES_SNAPSHOTS, snapshotDataAndSchemaFilename), fmt.Sprintf("%s", processor.config.HOST_SNAPSHOTS))
	// cmd5 := exec.Command("")

	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("failed to dump schema snapshot: %w", err)
	}

	if err := cmd2.Run(); err != nil {
		return fmt.Errorf("failed to dump data and schema snapshot: %w", err)
	}

	if err := cmd3.Run(); err != nil {
		return fmt.Errorf("failed to cp schema snapshot: %w", err)
	}

	if err := cmd4.Run(); err != nil {
		return fmt.Errorf("failed to cp data and schema snapshot: %w", err)
	}

	var fileContents [][]byte
	for _, file := range []string{snapshotDataAndSchemaFilename, snapshotSchemaFilename} {
		fileToAttach, err := os.Open(fmt.Sprintf("%s%s", processor.config.HOST_SNAPSHOTS, file))
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer fileToAttach.Close()

		fileContent, err := io.ReadAll(fileToAttach)
		if err != nil {
			return fmt.Errorf("failed to read file content: %w", err)
		}
		fileContents = append(fileContents, fileContent)
	}
	emailBody := fmt.Sprintf(`
		<h1>Hello Emilio Cliff</h1>
		<p>Your 3 days database snapshot</p>`)

	err := processor.sender.SendMail("Database Snapshot", emailBody, "text/plain", []string{"emiliocliff@gmail.com"}, nil, nil, []string{snapshotDataAndSchemaFilename, snapshotSchemaFilename}, fileContents)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// opts := []asynq.Option{
	// 	asynq.MaxRetry(2),
	// 	asynq.ProcessIn(72 * time.Hour),
	// 	asynq.Queue(QueueLow),
	// }

	// err = processor.distributor.DistributeTakeAndSendDBsnapshots(ctx, "word", opts...)
	// if err != nil {
	// 	return fmt.Errorf("Failed to schedule next snapshot: %w", err)
	// }

	log.Info().
		Str("type", task.Type()).
		Bytes("body", task.Payload()).
		Str("Success", "snapshot was succefuly sent to your email").
		Msg("tasked processed successfull")

	return nil
}
