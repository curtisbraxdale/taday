package handlers

import (
	"context"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/curtisbraxdale/taday/internal/auth"
	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    bool      `json:"priority"`
	RecurD      bool      `json:"recur_d"`
	RecurW      bool      `json:"recur_w"`
	RecurM      bool      `json:"recur_m"`
	RecurY      bool      `json:"recur_y"`
}

func (cfg *ApiConfig) GetUserEvents(w http.ResponseWriter, req *http.Request) {
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

	dbEvents, err := cfg.Queries.GetEventsByUserID(context.Background(), userID)
	if err != nil {
		log.Printf("Error finding events for given userID: %s", err)
		w.WriteHeader(500)
		return
	}
	events := []Event{}
	for _, e := range dbEvents {
		eventDescription := ""
		if e.Description.String != "" {
			eventDescription = e.Description.String
		}
		events = append(events, Event{ID: e.ID, UserID: e.UserID, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt, StartDate: e.StartDate, EndDate: e.EndDate, Title: e.Title, Description: eventDescription, Priority: e.Priority, RecurD: e.RecurD, RecurW: e.RecurW, RecurM: e.RecurM, RecurY: e.RecurY})
	}
	sortDir := req.URL.Query().Get("sort")
	if sortDir == "desc" {
		sort.Slice(events, func(i, j int) bool { return events[j].StartDate.Before(events[i].StartDate) })
	}
	respondWithJSON(w, 200, events)
}
