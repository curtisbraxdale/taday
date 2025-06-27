package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/curtisbraxdale/taday/internal/handlers"
	"github.com/curtisbraxdale/taday/internal/middleware"
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

	serveMux := http.NewServeMux()
	apiCfg := handlers.ApiConfig{Queries: dbQueries, Platform: platform, Secret: secret}

	serveMux.HandleFunc("GET /api/ready", handlers.ReadyCheck)
	serveMux.HandleFunc("POST /api/login", apiCfg.Login)
	serveMux.HandleFunc("POST /api/users", apiCfg.CreateUser)
	secure(serveMux, "POST /api/logout", apiCfg.Logout, secret)
	secure(serveMux, "POST /api/refresh", apiCfg.Refresh, secret)
	secure(serveMux, "POST /api/revoke", apiCfg.Revoke, secret)
	secure(serveMux, "POST /api/todos", apiCfg.CreateToDo, secret)
	secure(serveMux, "POST /api/events", apiCfg.CreateEvent, secret)
	secure(serveMux, "POST /api/tags", apiCfg.CreateTag, secret)
	secure(serveMux, "POST /api/events/{event_id}/tags", apiCfg.CreateEventTag, secret)
	secure(serveMux, "GET /api/users", apiCfg.GetUser, secret)
	secure(serveMux, "GET /api/events", apiCfg.GetUserEvents, secret)
	secure(serveMux, "GET /api/events/{event_id}", apiCfg.GetEvent, secret)
	secure(serveMux, "GET /api/tags", apiCfg.GetUserTags, secret)
	secure(serveMux, "GET /api/events/{event_id}/tags", apiCfg.GetEventTags, secret)
	secure(serveMux, "GET /api/todos", apiCfg.GetUserToDos, secret)
	secure(serveMux, "GET /api/todos/{todo_id}", apiCfg.GetToDo, secret)
	secure(serveMux, "PUT /api/users", apiCfg.UpdateUser, secret)
	secure(serveMux, "PUT /api/events/{event_id}", apiCfg.UpdateEvent, secret)
	secure(serveMux, "PUT /api/todos/{todo_id}", apiCfg.UpdateToDo, secret)

	server := http.Server{}
	server.Handler = serveMux
	server.Addr = ":8080"

	log.Fatal(server.ListenAndServe())
}

func secure(mux *http.ServeMux, methodAndPath string, handlerFunc http.HandlerFunc, secret string) {
	mux.Handle(methodAndPath, middleware.RequireAuth(secret, handlerFunc))
}
