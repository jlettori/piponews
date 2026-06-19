package models

import "time"

type User struct {
	ID                int64     `db:"id" json:"id"`
	Username          string    `db:"username" json:"username"`
	PasswordHash      string    `db:"password_hash" json:"-"`
	FirstName         string    `db:"first_name" json:"first_name"`
	LastName          string    `db:"last_name" json:"last_name"`
	Email             string    `db:"email" json:"email"`
	PreferredLanguage string    `db:"preferred_language" json:"preferred_language"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}

type Feed struct {
	ID            int64      `db:"id" json:"id"`
	UserID        int64      `db:"user_id" json:"user_id"`
	Title         string     `db:"title" json:"title"`
	URL           string     `db:"url" json:"url"`
	LastFetchedAt *time.Time `db:"last_fetched_at" json:"last_fetched_at"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	EntryCount    int        `db:"entry_count" json:"entry_count"`
}

type Entry struct {
	ID          int64      `db:"id" json:"id"`
	FeedID      int64      `db:"feed_id" json:"feed_id"`
	FeedTitle   string     `db:"-" json:"feed_title"`
	GUID        string     `db:"guid" json:"guid"`
	Title       string     `db:"title" json:"title"`
	URL         string     `db:"url" json:"url"`
	Summary     string     `db:"summary" json:"summary"`
	PublishedAt *time.Time `db:"published_at" json:"published_at"`
	IsSelected  bool       `db:"is_selected" json:"is_selected"`
}
