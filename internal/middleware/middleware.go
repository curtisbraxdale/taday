package middleware

import (
	"net/http"

	"github.com/curtisbraxdale/taday/internal/auth"
)

func RequireAuth(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized: No token", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateAccessToken(cookie.Value, secret)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := auth.ContextWithUserID(r.Context(), userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
