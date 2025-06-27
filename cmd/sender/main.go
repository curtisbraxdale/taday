package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/curtisbraxdale/taday/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	twilAccountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	twilAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilNumber := os.Getenv("TWILIO_PHONE_NUMBER")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
	}
	dbQueries := database.New(db)
	apiCfg := handlers.ApiConfig{Queries: dbQueries, Platform: platform, Secret: secret}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilAccountSid,
		Password: twilAuthToken,
	})

	day := time.Now().Weekday().String()
	userIDs, err := apiCfg.Queries.GetAllUsers(context.Background())
	if err != nil {
		log.Fatalf("Error getting userIDs: %s", err)
	}
	if day == "Monday" {
		for _, user := range userIDs {
			agenda, err := apiCfg.CreateWeeklyAgenda(user.ID)
			if err != nil {
				log.Printf("Error creating agenda for user %v : %s", user.Username, err)
			}
			params := &twilioApi.CreateMessageParams{}
			params.SetTo(user.PhoneNumber)
			params.SetFrom(twilNumber)
			params.SetBody(agenda)

			resp, err := client.Api.CreateMessage(params)
			if err != nil {
				fmt.Println("Error sending SMS message: " + err.Error())
			} else {
				response, _ := json.Marshal(*resp)
				fmt.Println("Response: " + string(response))
			}
		}
	} else {
		for _, user := range userIDs {
			agenda, err := apiCfg.CreateDailyAgenda(user.ID)
			if err != nil {
				log.Printf("Error creating agenda for user %v : %s", user.Username, err)
			}
			params := &twilioApi.CreateMessageParams{}
			params.SetTo(user.PhoneNumber)
			params.SetFrom(twilNumber)
			params.SetBody(agenda)

			resp, err := client.Api.CreateMessage(params)
			if err != nil {
				fmt.Println("Error sending SMS message: " + err.Error())
			} else {
				response, _ := json.Marshal(*resp)
				fmt.Println("Response: " + string(response))
			}
		}
	}
}
