package migrations

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbUrl, migrationsPath string) error {
	m, err := migrate.New("file://"+migrationsPath, dbUrl)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return srcErr
	}

	if dbErr != nil {
		return dbErr
	}

	return nil
}
