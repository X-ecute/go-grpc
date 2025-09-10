package db

import (
	"fmt"
	"os"

	"github.com/X-ecute/go-grpc/internal/rocket"
	"github.com/jmoiron/sqlx"
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

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return Store{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return Store{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return Store{db: db}, nil
}

func (s *Store) GetRocketByID(id string) (rocket.Rocket, error) {
	return rocket.Rocket{}, nil
}
func (s *Store) InsertRocket(rocket rocket.Rocket) (rocket.Rocket, error) {
	return rocket, nil
}

func (s *Store) DeleteRocket(id string) error {
	return nil
}
