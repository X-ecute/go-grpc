package db

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func (s *Store) Migrate() error {
	driver, err := postgres.WithInstance(s.db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migartions",
		"potgres",
		driver,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if !errors.Is(err, migrate.ErrNoChange) && err != nil {
		return err
	}
	return nil

}
