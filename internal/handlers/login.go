package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/curtisbraxdale/taday/internal/auth"
	"github.com/curtisbraxdale/taday/internal/database"
)

func (cfg *ApiConfig) Login(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	w.Header().Set("Access-Control-Allow-Origin", "https://taday.io")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	dbUser, err := cfg.Queries.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		log.Print("Incorrect email or password")
		w.WriteHeader(401)
		return
	}
	err = auth.CheckPasswordHash(dbUser.HashedPassword, params.Password)
	if err != nil {
		log.Print("Incorrect email or password")
		w.WriteHeader(401)
		return
	}
	// Create Access token.
	accessToken, err := auth.MakeAccessToken(dbUser.ID, cfg.Secret, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		w.WriteHeader(500)
	}
	// Create refresh token.
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		w.WriteHeader(500)
	}
	// Store refresh token in database.
	refTokenParams := database.CreateRefreshTokenParams{Token: refreshToken, ExpiresAt: sql.NullTime{Time: time.Now().Add(time.Hour * 24 * 60), Valid: true}, UserID: dbUser.ID, RevokedAt: sql.NullTime{Valid: false}}
	_, err = cfg.Queries.CreateRefreshToken(context.Background(), refTokenParams)
	if err != nil {
		log.Printf("Error storing refresh token: %s", err)
		w.WriteHeader(500)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(60 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	user := User{ID: dbUser.ID, CreatedAt: dbUser.CreatedAt, UpdatedAt: dbUser.UpdatedAt, Email: dbUser.Email}
	respondWithJSON(w, 200, user)
}
