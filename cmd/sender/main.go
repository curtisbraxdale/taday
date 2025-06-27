package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/curtisbraxdale/taday/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
	}
	dbQueries := database.New(db)
	apiCfg := handlers.ApiConfig{Queries: dbQueries, Platform: platform, Secret: secret}

	day := time.Now().Weekday().String()
	userIDs, err := apiCfg.Queries.GetAllUserIDs(context.Background())
	if err != nil {
		log.Fatalf("Error getting userIDs: %s", err)
	}
	if day == "Monday" {
		for _, id := range userIDs {
			agenda, err := apiCfg.CreateWeeklyAgenda(id)
			if err != nil {
				log.Printf("Error creating agenda for user %v : %s", id, err)
			}
			// TWILIO API SEND TEXT
		}
	} else {
		for _, id := range userIDs {
			agenda, err := apiCfg.CreateDailyAgenda(id)
			if err != nil {
				log.Printf("Error creating agenda for user %v : %s", id, err)
			}
			// TWILIO API SEND TEXT
		}
	}
}
