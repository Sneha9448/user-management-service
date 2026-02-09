package graph

import (
	"log"
	"time"
)

type Resolver struct{}

func (r *Resolver) TrackExecutionTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("Pure backend execution time for %s: %s", name, elapsed)
}
