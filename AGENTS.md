# piponews — RSS feed reader

## Tech stack

- **Language:** Go 1.26
- **Database:** SQLite (via `modernc.org/sqlite`, pure Go driver, no CGo)
- **DB query helpers:** `jmoiron/sqlx`
- **Web framework:** stdlib `net/http` with Go 1.22+ enhanced ServeMux (method-based routing)
- **Templating:** Go `html/template`, embedded via `embed.FS`
- **Frontend:** [DataStar](https://data-star.dev/) — reactive hypermedia (client-side signals, SSE-based DOM patching)
- **Auth:** session-based (signed cookies, `golang.org/x/crypto`)
- **RSS parsing:** `mmcdole/gofeed`
- **i18n:** custom internal package (en/fr/it)

## Project layout

```
main.go              — entry point, flag parsing, server start
routes.go            — HTTP route definitions
internal/
  db/sqlite.go       — DB init, goose-based schema migration
  db/migrations/     — SQL migration files (goose)
  models/models.go   — User, Feed, Entry structs
  handlers/
    entries.go       — entry listing/filtering/selection
    feeds.go         — feed CRUD, refresh, filter bar
    auth.go          — login/register/logout
    helpers.go       — shared utilities, sanitization, version
  templates/         — *.html templates
  i18n/              — keys.go + en/fr/it translations
  auth/auth.go       — session tokens, password hashing
static/              — CSS, client-side JS (datastar.js)
```

## Key conventions

- **Filter state** is stored client-side in DataStar signals (`feedFilter`, `dateFrom`, `dateTo`, `searchQuery`). Changes trigger `@get('/entries')` which sends signals as JSON via `datastar` query param.
- **Infinite scroll:** `/entries/more` fetches next page with offset.
- **Handlers** parse signals via `json.Unmarshal` from `r.URL.Query().Get("datastar")`.
- **Templates** use the `T` func for i18n and `dict` for passing multiple values.
- **Build version** is injected at build time via `-ldflags` into `handlers.BuildVersion`.

## Commands

```
make build   — go build -o piponews .
make run     — go run .
```

Server listens on `127.0.0.1:8080` by default (override with `-addr`).

## Migrations

SQL migrations are managed by [pressly/goose](https://github.com/pressly/goose) and live in `internal/db/migrations/`. Migrations run automatically on startup via `goose.Up`.

Create a new migration:
```
make migrate/new <name>
```

Other goose commands:
```
make migrate/status
make migrate/up
make migrate/down
```

## Tests

Standard Go test framework: `go test ./...`
