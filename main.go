package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jlettori/piponews/internal/auth"
	"github.com/jlettori/piponews/internal/db"
	"github.com/jlettori/piponews/internal/handlers"
	"github.com/jlettori/piponews/internal/i18n"
)

//go:embed static
var staticAssets embed.FS

func staticFileSystem() http.FileSystem {
	sub, err := fs.Sub(staticAssets, "static")
	if err != nil {
		log.Fatalf("failed to create static file system: %v", err)
	}
	return http.FS(sub)
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	versionFlag := flag.Bool("version", false, "print version and exit")
	verbose := flag.Bool("verbose", false, "log startup steps")
	flag.Parse()

	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	v := func(format string, args ...any) {
		if *verbose {
			log.Printf("▶ "+format, args...)
		}
	}

	v("resolving database path")
	appName := i18n.Global.T(i18n.En, i18n.AppName)
	dbPath, pathErr := filepath.Abs(appName + ".db")
	if pathErr != nil {
		log.Fatalf("failed to resolve database path: %v", pathErr)
	}
	v("database path: %s", dbPath)

	v("initializing database")
	db.SetVerbose(*verbose)
	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	v("database initialized")

	v("creating session store")
	sessions := auth.NewStore(database)
	v("building router")
	mux := newRouter(database, sessions)

	log.Printf("Init\n%s %s (build %s)\ndatabase: %s\nlistening on %s",
		appName, Version, handlers.BuildVersion, dbPath, *addr)
	v("starting HTTP server on %s", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
