# piponews

[![Build](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/jlettori/piponews)
[![Coverage](https://img.shields.io/badge/coverage-76%25-25A162)](https://github.com/jlettori/piponews)
[![Go Reference](https://img.shields.io/badge/docs-pkg.go.dev-007D9C)](https://pkg.go.dev/github.com/jlettori/piponews)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)

Multi-user RSS feed reader with a hypermedia frontend.

## Installation

```bash
# Build from source
make build
```

Or run directly during development:

```bash
make run
```

## Usage

Start the server and open http://127.0.0.1:8080:

```bash
./piponews
# → piponews 0.1.0
# → database: ./piponews.db
# → listening on 127.0.0.1:8080
```

Register an account, add a feed URL, and refresh.

## Flags

```
-addr string      listen address (default "127.0.0.1:8080")
--version         print version and exit
```
