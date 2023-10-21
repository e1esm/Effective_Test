package migrations

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type MigrationType int

const (
	maxRetries               = 10
	UP         MigrationType = iota
	DOWN
	fileDir = "file:///migrations"
)

func ConnectAndRunMigrations(ctx context.Context, dbURL string, migrationType MigrationType) (*pgxpool.Pool, error) {
	pool, err := connect(ctx, dbURL)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if err := runMigrations(dbURL, fileDir, migrationType); err != nil {
		return nil, err
	}
	return pool, nil
}

func connect(ctx context.Context, dbURl string) (*pgxpool.Pool, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	var pool *pgxpool.Pool
	var err error
	retry := 0
	for retry < maxRetries {
		select {
		case <-ticker.C:
			pool, err = pgxpool.New(ctx, dbURl)
			retry++
			if err := pool.Ping(context.Background()); err == nil {
				return pool, nil
			}
		}
	}
	return nil, fmt.Errorf("coudln't have connected to the database: %v", err)
}

func runMigrations(dbURL, fileDir string, migrationType MigrationType) error {
	m, err := migrate.New(fileDir, dbURL)
	if err != nil {
		return err
	}
	switch migrationType {
	case UP:
		if err = m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				return nil
			}
			log.Println(err.Error())
			return err
		}
	case DOWN:
		if err = m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				return nil
			}
			log.Println(err.Error())
			return err
		}
	}
	return nil
}
