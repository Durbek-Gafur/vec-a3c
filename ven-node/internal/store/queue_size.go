package store

import (
	"context"
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)



func (s *MySQLStore) GetQueueSize(ctx context.Context) (int, error) {
	// var size int
	// err := s.db.QueryRowContext(ctx, "SELECT size FROM queue_size WHERE id = 1").Scan(&size)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return 0, &StoreError{Err: ErrNotFound, StatusCode: http.StatusNotFound}
	// 	}
	// 	return 0, err
	// }

	return s.GetQueueSizeFromDBorENV(ctx)
}

func (s *MySQLStore) SetQueueSize(ctx context.Context, size int) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO queue_size (id, size) VALUES (1, ?) ON DUPLICATE KEY UPDATE size = ?", size, size)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) UpdateQueueSize(ctx context.Context, size int) error {
	_, err := s.db.ExecContext(ctx, "UPDATE queue_size SET size = ? WHERE id = 1", size)
	if err != nil {
		return err
	}
	return nil
}

func (s *MySQLStore) getQueueSizeFromDB(ctx context.Context) (int, error) {
	
	queueSizeQuery := "SELECT size FROM queue_size LIMIT 1;"
	var queueSize int
	err := s.db.QueryRowContext(ctx, queueSizeQuery).Scan(&queueSize)
	if err != nil {
		return 0, err
	}
	// fmt.Printf("queue_size %d from DB\n", queueSize)
	return queueSize, nil
}

func (s *MySQLStore) getQueueSizeFromEnv() (int, error) {
	queueSizeEnv := os.Getenv("QUEUE_SIZE")
	if queueSizeEnv == "" {
		fmt.Println("QUEUE_SIZE not set, using default value 10")
		return 10, nil //errors.New("QUEUE_SIZE not set, using default value 10")
	}
	queueSize, err := strconv.Atoi(queueSizeEnv)
	if err != nil {
		return 0, fmt.Errorf("invalid value for QUEUE_SIZE: %s", queueSizeEnv)
	}
	// fmt.Printf("QUEUE_SIZE from ENV %d", queueSize)
	return queueSize, nil
}

func (s *MySQLStore) GetQueueSizeFromDBorENV(ctx context.Context) (int, error) {
	queueSize, err := s.getQueueSizeFromDB(ctx)
	if err == nil {
		return queueSize, nil
	}

	queueSize, err = s.getQueueSizeFromEnv()
	if err != nil {
		return 0, err
	}

	return queueSize, nil
}