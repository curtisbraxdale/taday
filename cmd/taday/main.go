package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/curtisbraxdale/taday/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATBASE_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
	}
	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	apiCfg := handlers.ApiConfig{Queries: dbQueries, Platform: platform, Secret: secret}

	serveMux.HandleFunc("GET /api/ready", handlers.Ready)
	serveMux.HandleFunc("GET /api/users", apiCfg.GetUsers)
	serveMux.HandleFunc("GET /api/events/{userID}", apiCfg.GetUserEvents)

	server := http.Server{}
	server.Handler = serveMux
	server.Addr = ":8080"

	log.Fatal(server.ListenAndServe())
}
