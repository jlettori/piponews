# piponews

Multi-user RSS feed reader with [DataStar](https://data-star.dev/) hypermedia frontend. Follow feeds, browse entries with filters, select items, and export as HTML or plain text.

## Quick start

```bash
make build
./piponews
# → piponews 0.1.0
# → database: /home/user/piponews/piponews.db
# → listening on 127.0.0.1:8080
```

The build embeds a timestamp in the binary which is appended to static asset URLs (`?v=<timestamp>`) to prevent stale browser caches. Use `make run` during development.
```bash
make run
```

Open http://127.0.0.1:8080, register an account, add a feed URL, and refresh.

## Features

- **Multi-user** — register/login with bcrypt + HMAC-signed session cookies
- **Feed management** — add/remove feeds, refresh individually or all at once, alphabetically sorted
- **Entry browsing** — most-recent-first, filter by feed or date range, auto-filter on feed select
- **Selection & export** — select entries for export as HTML or plain text, opens in new tab
- **Internationalization** — auto-detects browser language (English / French)
- **DataStar UI** — reactive frontend with server-driven DOM morphing, no JS build step

## Flags

```
-addr string      listen address (default "127.0.0.1:8080")
--version         print version and exit
```

## Stack

- **Backend** — Go 1.26, `net/http` with pattern routing
- **Frontend** — DataStar (hypermedia framework, `data-*` attributes, SSE)
- **Database** — SQLite via `modernc.org/sqlite` (pure Go, no CGo)
- **RSS** — `github.com/mmcdole/gofeed`

## Project structure

```
piponews/
├── main.go                    # entry point, flag parsing
├── routes.go                  # route registration
├── version.go                 # semver version constant
├── internal/
│   ├── auth/auth.go           # bcrypt, session cookies, middleware
│   ├── db/sqlite.go           # SQLite init, schema migrations
│   ├── handlers/
│   │   ├── auth.go            # login, register, logout
│   │   ├── feeds.go           # list, create, delete, refresh
│   │   ├── entries.go         # list, filter, toggle select
│   │   ├── export.go          # HTML/text file generation
│   │   └── helpers.go         # shared utilities
│   ├── i18n/                  # translations
│   │   ├── i18n.go            # engine, locale detection
│   │   ├── keys.go            # translatable string keys
│   │   ├── en.go              # English
│   │   └── fr.go              # French
│   ├── models/models.go       # User, Feed, Entry structs
│   └── templates/             # Go html/template + DataStar attributes
├── static/
│   ├── datastar.js            # self-hosted DataStar bundle
│   └── style.css              # application styles
├── go.mod / go.sum
```

## API routes

| Method | Route | Description |
|--------|-------|-------------|
| GET/POST | `/login` | Sign in |
| GET/POST | `/register` | Create account |
| POST | `/logout` | Sign out |
| GET | `/feeds` | Main app page |
| POST | `/feeds` | Add a feed (fetches entries immediately) |
| DELETE | `/feeds/{id}` | Remove a feed |
| POST | `/feeds/{id}/refresh` | Refresh one feed |
| POST | `/entries/refresh-all` | Refresh all feeds |
| GET | `/entries` | List entries (filtered by `datastar` signals) |
| POST | `/entries/{id}/toggle-select` | Toggle export selection |
| POST | `/entries/select-all` | Select all visible entries |
| POST | `/entries/clear-selection` | Clear all selections |
| POST | `/entries/export` | Open selected entries as HTML or txt |
