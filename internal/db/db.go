package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/X-ecute/go-grpc/internal/rocket"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Store struct {
	db *sqlx.DB
}

// NewStore - return a new store or error
func New() (Store, error) {
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDatabase := os.Getenv("POSTGRES_DATABASE")
	postgresSSLMode := os.Getenv("POSTGRES_SSLMODE")
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		postgresUser,
		postgresPassword,
		postgresHost,
		postgresPort,
		postgresDatabase,
		postgresSSLMode,
	)

	var db *sqlx.DB
	var err error

	// Retry connection for 30 seconds
	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("postgres", connectionString)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/10): %v", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return Store{}, fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return Store{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return Store{db: db}, nil
}
func (s *Store) GetRocketByID(ctx context.Context, id string) (rocket.Rocket, error) {
	var rkt rocket.Rocket

	query := `SELECT id, name, type FROM rockets WHERE id = $1`

	err := s.db.GetContext(ctx, &rkt, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return rocket.Rocket{}, status.Errorf(codes.NotFound, "rocket with id %s not found", id)
		}
		return rocket.Rocket{}, status.Errorf(codes.Internal, "failed to get rocket: %v", err)
	}

	return rkt, nil
}

func (s *Store) InsertRocket(ctx context.Context, rkt rocket.Rocket) (rocket.Rocket, error) {
	_, err := s.db.NamedExecContext(ctx, `INSERT INTO rockets (id, name, type) VALUES (:id, :name, :type)`, rkt)
	if err != nil {
		return rocket.Rocket{}, err
	}

	return rocket.Rocket{
		ID:   rkt.ID,
		Name: rkt.Name,
		Type: rkt.Type,
	}, nil
}

func (s *Store) DeleteRocket(ctx context.Context, id string) error {
	query := `DELETE FROM rockets WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rocket: %w", err)
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rocket with id %s not found", id)
	}

	return nil
}
