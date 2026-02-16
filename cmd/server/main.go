package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"user-management-service/graph"
	"user-management-service/internal/config"
	"user-management-service/internal/database"
	"user-management-service/internal/email"
	"user-management-service/internal/middleware"
	"user-management-service/internal/router"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()
	log.Println("Starting User Management Service...")

	// 2. Initialize Email Service
	email.Init(cfg)

	// 3. Connect to Database (blocks until connected or fails)
	database.ConnectDB(cfg.DatabaseURL)
	defer database.CloseDB()

	// 3. Setup Router
	r := router.SetupRouter()

	// GraphQL Handler
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Config: cfg}}))
	r.Handle("/graphql", srv)
	r.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	// 4. Start Server
	log.Printf("Server listening on port %s", cfg.Port)

	// Start pprof server in a goroutine
	go func() {
		log.Println("pprof running on http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Printf("pprof server failed: %v", err)
		}
	}()

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := middleware.AuthMiddleware()(r)
	handler = c.Handler(handler)

	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
