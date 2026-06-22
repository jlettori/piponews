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
	flag.Parse()

	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	appName := i18n.Global.T(i18n.En, i18n.AppName)
	dbPath, pathErr := filepath.Abs(appName + ".db")
	if pathErr != nil {
		log.Fatalf("failed to resolve database path: %v", pathErr)
	}
	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	sessions := auth.NewStore(database)
	mux := newRouter(database, sessions)

	log.Printf("Init\n%s %s (build %s)\ndatabase: %s\nlistening on %s",
		appName, Version, handlers.BuildVersion, dbPath, *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
