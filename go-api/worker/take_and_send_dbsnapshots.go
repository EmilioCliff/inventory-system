package worker

import (
	"compress/gzip"
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
		return fmt.Errorf("failed to enqueue task take and send dbsnapshots: %w", err)
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

	snapshotSchemaFilename := timestamp + "_schema_snapshot.tar"

	cmd := exec.Command(
		"pg_dump",
		"-U", os.Getenv("PGUSER"),
		"-h", os.Getenv("PGHOST"),
		"-p", os.Getenv("PGPORT"),
		"-d", os.Getenv("POSTGRES_DB"),
		"-F", "t",
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", os.Getenv("POSTGRES_PASSWORD")))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to run pg_dump: %v\nOutput: %s", err, output)
	}

	cmd.Env = os.Environ()

	if err := os.WriteFile(snapshotSchemaFilename, output, 0644); err != nil {
		return fmt.Errorf("Failed to write snapshot to file: %v", err)
	}

	compressedSnapshotFilename := snapshotSchemaFilename + ".gz"
	err = compressFile(snapshotSchemaFilename, compressedSnapshotFilename)
	if err != nil {
		return fmt.Errorf("Failed to compress snapshot file: %v", err)
	}

	fileContent, err := os.ReadFile(compressedSnapshotFilename)
	if err != nil {
		return fmt.Errorf("Failed to read file content: %w", err)
	}

	emailBody := fmt.Sprintf(`
		<h1>Hello Emilio Cliff</h1>
		<p>Your 3 days database snapshot</p>`)

	err = processor.sender.SendMail(
		"Database Snapshot",
		emailBody,
		"text/plain",
		[]string{"clifftest33@gmail.com"},
		nil,
		nil,
		[]string{snapshotSchemaFilename},
		[][]byte{fileContent},
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if err := os.Remove(snapshotSchemaFilename); err != nil {
		return fmt.Errorf("Failed to delete snapshot file: %w", err)
	}

	if err := os.Remove(compressedSnapshotFilename); err != nil {
		return fmt.Errorf("Failed to delete compressed snapshot file: %w", err)
	}

	// Uncomment and adjust as needed
	// opts := []asynq.Option{
	// 	asynq.MaxRetry(2),
	// 	asynq.ProcessIn(24 * time.Hour),
	// 	asynq.Queue(QueueLow),
	// }
	// err = processor.distributor.DistributeTakeAndSendDBsnapshots(ctx, "word", opts...)
	// if err != nil {
	// 	return fmt.Errorf("Failed to schedule next snapshot: %w", err)
	// }

	log.Info().
		Str("type", task.Type()).
		Str("Success", "snapshot was successfully sent to your email").
		Msg("task processed successfully")

	return nil
}

func compressFile(source, target string) error {
	inFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, inFile)
	return err
}
