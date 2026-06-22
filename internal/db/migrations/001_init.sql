-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    preferred_language TEXT NOT NULL DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL,
    last_fetched_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, url)
);

CREATE TABLE IF NOT EXISTS entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id INTEGER NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    guid TEXT NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL DEFAULT '',
    published_at DATETIME,
    UNIQUE(feed_id, guid)
);

CREATE TABLE IF NOT EXISTS entry_selections (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entry_id INTEGER NOT NULL REFERENCES entries(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, entry_id)
);

CREATE INDEX IF NOT EXISTS idx_entries_feed_id ON entries(feed_id);
CREATE INDEX IF NOT EXISTS idx_entries_published_at ON entries(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_entries_feed_pub ON entries(feed_id, published_at DESC);

-- +goose Down
DROP TABLE IF EXISTS entry_selections;
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS feeds;
DROP TABLE IF EXISTS users;
