package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sort"
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

func (cfg *ApiConfig) GetUserToDos(w http.ResponseWriter, req *http.Request) {
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

	dbToDos, err := cfg.Queries.GetTodosByUserID(req.Context(), userID)
	if err != nil {
		log.Printf("Error finding events for given userID: %s", err)
		w.WriteHeader(500)
		return
	}
	toDos := []ToDo{}
	for _, t := range dbToDos {
		toDoDescription := ""
		if t.Description.String != "" {
			toDoDescription = t.Description.String
		}
		toDos = append(toDos, ToDo{ID: t.ID, UserID: t.UserID, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt, Date: t.Date.Time, Title: t.Title, Description: toDoDescription})
	}
	sortDir := req.URL.Query().Get("sort")
	if sortDir == "desc" {
		sort.Slice(toDos, func(i, j int) bool { return toDos[j].Date.Before(toDos[i].Date) })
	}
	respondWithJSON(w, 200, toDos)
}

func (cfg *ApiConfig) GetToDo(w http.ResponseWriter, req *http.Request) {
	toDoID, err := uuid.Parse(req.PathValue("todo_id"))
	if err != nil {
		log.Printf("Error parsing uuid: %s", err)
		w.WriteHeader(500)
		return
	}
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

	dbToDo, err := cfg.Queries.GetTodoByID(req.Context(), toDoID)
	if err != nil {
		log.Printf("Error finding events for given userID: %s", err)
		w.WriteHeader(500)
		return
	}
	if dbToDo.UserID != userID {
		log.Printf("Unauthorized access: user %s tried to access todo %s owned by %s", userID, dbToDo.ID, dbToDo.UserID)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	toDo := ToDo{ID: dbToDo.ID, UserID: dbToDo.UserID, CreatedAt: dbToDo.CreatedAt, UpdatedAt: dbToDo.UpdatedAt, Date: dbToDo.Date.Time, Title: dbToDo.Title, Description: dbToDo.Description.String}
	respondWithJSON(w, 200, toDo)
}

func (cfg *ApiConfig) UpdateToDo(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Date        time.Time `json:"date"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}

	toDoID, err := uuid.Parse(req.PathValue("todo_id"))
	if err != nil {
		log.Printf("Error parsing uuid: %s", err)
		w.WriteHeader(500)
		return
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

	_, err = auth.ValidateAccessToken(accessToken, cfg.Secret)
	if err != nil {
		log.Printf("Access token invalid: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	dbTodoParams := database.UpdateToDoParams{Date: sql.NullTime{Time: params.Date, Valid: true}, Title: params.Title, Description: sql.NullString{String: params.Description, Valid: true}, TodoID: toDoID}
	dbTodo, err := cfg.Queries.UpdateToDo(req.Context(), dbTodoParams)
	if err != nil {
		log.Printf("Error updating todo: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	toDo := ToDo{ID: dbTodo.ID, UserID: dbTodo.UserID, CreatedAt: dbTodo.CreatedAt, UpdatedAt: dbTodo.UpdatedAt, Date: dbTodo.Date.Time, Title: dbTodo.Title, Description: dbTodo.Description.String}
	respondWithJSON(w, 201, toDo)
}
