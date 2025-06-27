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

func (cfg *ApiConfig) CreateEvent(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
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
	dbEventParams := database.CreateEventParams{UserID: userID, StartDate: params.StartDate, EndDate: params.EndDate, Title: params.Title, Description: sql.NullString{String: params.Description, Valid: true}, Priority: params.Priority, RecurD: params.RecurD, RecurW: params.RecurW, RecurM: params.RecurM, RecurY: params.RecurY}
	dbEvent, err := cfg.Queries.CreateEvent(req.Context(), dbEventParams)
	if err != nil {
		log.Printf("Error creating todo: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	event := Event{ID: dbEvent.ID, UserID: dbEvent.UserID, CreatedAt: dbEvent.CreatedAt, UpdatedAt: dbEvent.UpdatedAt, StartDate: dbEvent.StartDate, EndDate: dbEvent.EndDate, Title: dbEvent.Title, Description: dbEvent.Description.String, Priority: dbEvent.Priority, RecurD: dbEvent.RecurD, RecurW: dbEvent.RecurW, RecurM: dbEvent.RecurM, RecurY: dbEvent.RecurY}
	respondWithJSON(w, 201, event)
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

	tagFilter := req.URL.Query().Get("tag")
	rangeFilter := req.URL.Query().Get("range")
	now := time.Now()

	var startDate, endDate time.Time

	switch rangeFilter {
	case "day":
		startDate = now.Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	case "week":
		weekday := int(now.Weekday())
		startDate = now.AddDate(0, 0, -weekday)
		endDate = startDate.AddDate(0, 0, 7)
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)
	case "year":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(1, 0, 0)
	default:
		startDate = time.Time{}
		endDate = time.Time{}
	}

	dbEventParams := database.GetFilteredEventsParams{UserID: userID, StartDate: startDate, EndDate: endDate, Tag: tagFilter}
	dbEvents, err := cfg.Queries.GetFilteredEvents(req.Context(), dbEventParams)
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
	/*
		 * Old non-filtered query
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
	*/
}

func (cfg *ApiConfig) GetEvent(w http.ResponseWriter, req *http.Request) {
	eventID, err := uuid.Parse(req.PathValue("event_id"))
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

	dbEvent, err := cfg.Queries.GetEventByID(req.Context(), eventID)
	if err != nil {
		log.Printf("Error finding event for given id: %s", err)
		w.WriteHeader(500)
		return
	}
	if dbEvent.UserID != userID {
		log.Printf("Unauthorized access: user %s tried to access event %s owned by %s", userID, dbEvent.ID, dbEvent.UserID)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	event := Event{ID: dbEvent.ID, UserID: dbEvent.UserID, CreatedAt: dbEvent.CreatedAt, UpdatedAt: dbEvent.UpdatedAt, StartDate: dbEvent.StartDate, EndDate: dbEvent.EndDate, Title: dbEvent.Title, Description: dbEvent.Description.String, Priority: dbEvent.Priority, RecurD: dbEvent.RecurD, RecurW: dbEvent.RecurW, RecurM: dbEvent.RecurM, RecurY: dbEvent.RecurY}
	respondWithJSON(w, 200, event)
}

func (cfg *ApiConfig) UpdateEvent(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
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

	eventID, err := uuid.Parse(req.PathValue("event_id"))
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
	dbEventParams := database.UpdateEventParams{StartDate: params.StartDate, EndDate: params.EndDate, Title: params.Title, Description: sql.NullString{String: params.Description, Valid: true}, Priority: params.Priority, RecurD: params.RecurD, RecurW: params.RecurW, RecurM: params.RecurM, RecurY: params.RecurY, EventID: eventID}
	dbEvent, err := cfg.Queries.UpdateEvent(req.Context(), dbEventParams)
	if err != nil {
		log.Printf("Error updating event: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	event := Event{ID: dbEvent.ID, UserID: dbEvent.UserID, CreatedAt: dbEvent.CreatedAt, UpdatedAt: dbEvent.UpdatedAt, StartDate: dbEvent.StartDate, EndDate: dbEvent.EndDate, Title: dbEvent.Title, Description: dbEvent.Description.String, Priority: dbEvent.Priority, RecurD: dbEvent.RecurD, RecurW: dbEvent.RecurW, RecurM: dbEvent.RecurM, RecurY: dbEvent.RecurY}
	respondWithJSON(w, 201, event)
}

func (cfg *ApiConfig) DeleteEvent(w http.ResponseWriter, req *http.Request) {
	eventID, err := uuid.Parse(req.PathValue("event_id"))
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

	_, err = auth.ValidateAccessToken(accessToken, cfg.Secret)
	if err != nil {
		log.Printf("Access token invalid: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.Queries.DeleteEventByID(req.Context(), eventID)
	if err != nil {
		log.Printf("Error deleting event: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}
