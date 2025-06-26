package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userIDContextKey = contextKey("userID")

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return userID, ok
}

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func MakeAccessToken(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	unsignedAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "taday", IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)), Subject: userID.String()})
	signedAccessToken, err := unsignedAccessToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedAccessToken, nil
}

func ValidateAccessToken(tokenString, tokenSecret string) (uuid.UUID, error) {
	accessToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := accessToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, errors.New("Invalid token.")
	}
	subject, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func MakeRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)
	return token, nil
}
