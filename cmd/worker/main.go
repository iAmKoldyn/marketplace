package main

import (
	"log"
	"os"

	"github.com/iAmKoldyn/marketplace/internal/worker"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR") // e.g. "redis:6379"
	p := worker.NewProcessor(redisAddr)
	log.Println("worker: starting processor…")
	if err := p.Run(); err != nil {
		log.Fatalf("worker: processor error: %v", err)
	}
}
