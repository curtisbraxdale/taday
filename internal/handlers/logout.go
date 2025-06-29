package handlers

import (
	"net/http"
	"time"
)

func (cfg *ApiConfig) Logout(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://taday.io")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if refreshCookie, err := req.Cookie("refresh_token"); err == nil {
		_ = cfg.Queries.RevokeToken(req.Context(), refreshCookie.Value)
	}

	expire := func(name string) {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}
	expire("access_token")
	expire("refresh_token")

	w.WriteHeader(http.StatusOK)
}
