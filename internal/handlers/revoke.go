package handlers

import (
	"context"
	"log"
	"net/http"
)

func (cfg *ApiConfig) Revoke(w http.ResponseWriter, req *http.Request) {
	// Get refresh token from cookies.
	refreshCookie, err := req.Cookie("refresh_token")
	if err != nil {
		log.Printf("Refresh token not found in cookies: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	refreshToken := refreshCookie.Value
	err = cfg.Queries.RevokeToken(context.Background(), refreshToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err)
		w.WriteHeader(401)
		return
	}
	w.WriteHeader(204)
	return
}
