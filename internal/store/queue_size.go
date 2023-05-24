package store

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)



func (s *MySQLStore) GetQueueSize(ctx context.Context) (int, error) {
	var size int
	err := s.db.QueryRowContext(ctx, "SELECT size FROM queue_size WHERE id = 1").Scan(&size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, &StoreError{Err: ErrNotFound, StatusCode: http.StatusNotFound}
		}
		return 0, err
	}

	return size, nil
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

