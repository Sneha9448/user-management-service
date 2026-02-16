package graph

import (
	"log"
	"time"
	"user-management-service/internal/config"
)

type Resolver struct {
	Config *config.Config
}

func (r *Resolver) TrackExecutionTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("Pure backend execution time for %s: %s", name, elapsed)
}
