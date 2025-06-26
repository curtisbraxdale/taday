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
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
}

func (cfg *ApiConfig) GetUser(w http.ResponseWriter, req *http.Request) {
	accessCookie, err := req.Cookie("access_token")
	if err != nil {
		log.Printf("Access token not found in cookies: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	accessToken := accessCookie.Value

	userID, err := auth.ValidateAccessToken(accessToken, cfg.Secret)
	if err != nil {
		log.Printf("Access token invalid: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbUser, err := cfg.Queries.GetUserByID(context.Background(), userID)
	if err != nil {
		log.Printf("Error getting user from userID: %s", err)
		w.WriteHeader(500)
		return
	}
	user := User{ID: dbUser.ID, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.UpdatedAt, Username: dbUser.Username, Email: dbUser.Email, PhoneNumber: dbUser.PhoneNumber}
	respondWithJSON(w, 200, user)
}

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phone_number"`
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

	dbUserParams := database.CreateUserParams{Username: params.Username, Email: params.Email, HashedPassword: hashedPassword, PhoneNumber: params.PhoneNumber}
	dbUser, err := cfg.Queries.CreateUser(context.Background(), dbUserParams)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	newUser := User{ID: dbUser.ID, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.UpdatedAt, Username: dbUser.Username, Email: dbUser.Email, PhoneNumber: dbUser.PhoneNumber}
	respondWithJSON(w, 201, newUser)
}
