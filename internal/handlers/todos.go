package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/curtisbraxdale/taday/internal/auth"
	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/google/uuid"
)

type ToDo struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (cfg *ApiConfig) CreateToDo(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Date        time.Time `json:"date"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}
	accessCookie, err := req.Cookie("access_token")
	if err != nil {
		log.Printf("Access token not found in cookies: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	accessToken := accessCookie.Value

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	userID, err := auth.ValidateAccessToken(accessToken, cfg.Secret)
	if err != nil {
		log.Printf("Access token invalid: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	dbTodoParams := database.CreateTodoParams{UserID: userID, Date: sql.NullTime{Time: params.Date, Valid: true}, Title: params.Title, Description: sql.NullString{String: params.Description, Valid: true}}
	dbTodo, err := cfg.Queries.CreateTodo(req.Context(), dbTodoParams)
	if err != nil {
		log.Printf("Error creating todo: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	toDo := ToDo{ID: dbTodo.ID, UserID: dbTodo.UserID, CreatedAt: dbTodo.CreatedAt, UpdatedAt: dbTodo.UpdatedAt, Date: dbTodo.Date.Time, Title: dbTodo.Title, Description: dbTodo.Description.String}
	respondWithJSON(w, 201, toDo)
}
