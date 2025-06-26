package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/curtisbraxdale/taday/internal/auth"
)

func (cfg *ApiConfig) Refresh(w http.ResponseWriter, req *http.Request) {
	refreshCookie, err := req.Cookie("refresh_token")
	if err != nil {
		log.Printf("Refresh token not found in cookies: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	refreshToken := refreshCookie.Value

	// Get associated user from database.
	dbRefToken, err := cfg.Queries.GetUserByToken(req.Context(), refreshToken)
	if err != nil {
		log.Printf("Invalid refresh token: %s", err)
		w.WriteHeader(401)
		return
	}
	// Check if token has been revoked.
	if dbRefToken.RevokedAt.Valid || !dbRefToken.ExpiresAt.Valid || dbRefToken.ExpiresAt.Time.Before(time.Now()) {
		w.WriteHeader(401)
		return
	}
	// Create new access token that expires in 1 hour.
	token := ""
	token, err = auth.MakeAccessToken(dbRefToken.UserID, cfg.Secret, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		w.WriteHeader(500)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(200)
}
