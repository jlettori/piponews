package db

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	*sqlx.DB
}

func InitDB(path string) (*DB, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &DB{DB: db}, nil
}

func migrate(db *sqlx.DB) error {
	goose.SetBaseFS(migrationsFS)
	goose.SetDialect("sqlite3")
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
