package main

import (
	"log"
	"net/http"

	"github.com/yourusername/marketplace/internal/config"
	"github.com/yourusername/marketplace/internal/httpapi"
	"github.com/yourusername/marketplace/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	db, err := store.Connect(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()

	if err := store.RunMigrations(cfg.PostgresURL); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	router := httpapi.NewRouter(store.NewQueries(db), cfg)
	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
