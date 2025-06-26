package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/curtisbraxdale/taday/internal/auth"
	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *ApiConfig) GetUser(w http.ResponseWriter, req *http.Request) {
	userID, err := uuid.Parse(req.PathValue("userID"))
	if err != nil {
		log.Printf("Error parsing uuid: %s", err)
		w.WriteHeader(500)
		return
	}
	dbUser, err := cfg.Queries.GetUserByID(context.Background(), userID)
	if err != nil {
		log.Printf("Error getting user from userID: %s", err)
		w.WriteHeader(500)
		return
	}
	user := User{ID: dbUser.ID, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.UpdatedAt, Username: dbUser.Username, Email: dbUser.Email}
	respondWithJSON(w, 200, user)
}

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, req http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}

	dbUserParams := database.CreateUserParams{Username: params.Username, Email: params.Email, HashedPassword: hashedPassword}
	dbUser, err := cfg.Queries.CreateUser(context.Background(), dbUserParams)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	newUser := User{ID: dbUser.ID, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.UpdatedAt, Username: dbUser.Username, Email: dbUser.Email}
	respondWithJSON(w, 201, newUser)
}
