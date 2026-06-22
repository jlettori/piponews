package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

var verbose bool

func SetVerbose(v bool) {
	verbose = v
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	*sqlx.DB
}

func InitDB(path string) (*DB, error) {
	if verbose {
		log.Printf("  opening SQLite database at %s", path)
	}
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if verbose {
		log.Printf("  database opened, setting connection limits")
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if verbose {
		log.Printf("  running migrations")
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	if verbose {
		log.Printf("  migrations complete")
	}

	return &DB{DB: db}, nil
}

func migrate(db *sqlx.DB) error {
	if verbose {
		log.Printf("  configuring goose")
	}
	goose.SetBaseFS(migrationsFS)
	goose.SetDialect("sqlite3")
	if verbose {
		log.Printf("  running goose.Up")
	}
	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}
	return nil
}

func NullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func NullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
