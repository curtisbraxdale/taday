package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/curtisbraxdale/taday/internal/auth"
	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/google/uuid"
)

type Tag struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Color  string    `json:"color"`
}

type EventTag struct {
	EventID uuid.UUID `json:"event_id"`
	TagID   uuid.UUID `json:"tag_id"`
}

func (cfg *ApiConfig) CreateTag(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Name  string `json:"name"`
		Color string `json:"color"`
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

	dbTagParams := database.CreateTagParams{UserID: userID, Name: params.Name, Color: params.Color}
	dbTag, err := cfg.Queries.CreateTag(req.Context(), dbTagParams)
	if err != nil {
		log.Printf("Error creating tag: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tag := Tag{ID: dbTag.ID, UserID: dbTag.UserID, Name: dbTag.Name, Color: dbTag.Color}
	respondWithJSON(w, 201, tag)
}

func (cfg *ApiConfig) CreateEventTag(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		TagID uuid.UUID `json:"tag_id"`
	}
	eventID, err := uuid.Parse(req.PathValue("event_id"))
	if err != nil {
		log.Printf("Invalid event_id in path: %s", err)
		w.WriteHeader(http.StatusBadRequest)
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

	dbEventTagParams := database.CreateEventTagParams{EventID: eventID, TagID: params.TagID}
	dbEventTag, err := cfg.Queries.CreateEventTag(req.Context(), dbEventTagParams)
	if err != nil {
		log.Printf("Error creating tag: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventTag := EventTag{EventID: dbEventTag.EventID, TagID: dbEventTag.TagID}
	respondWithJSON(w, 201, eventTag)
}

func (cfg *ApiConfig) GetUserTags(w http.ResponseWriter, req *http.Request) {
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

	dbTags, err := cfg.Queries.GetTagsByUserID(req.Context(), userID)
	if err != nil {
		log.Printf("Error finding tags for given userID: %s", err)
		w.WriteHeader(500)
		return
	}
	tags := []Tag{}
	for _, t := range dbTags {
		tags = append(tags, Tag{ID: t.ID, UserID: t.UserID, Name: t.Name, Color: t.Color})
	}
	respondWithJSON(w, 200, tags)
}

func (cfg *ApiConfig) GetEventTags(w http.ResponseWriter, req *http.Request) {
	eventID, err := uuid.Parse(req.PathValue("id"))
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

	dbTags, err := cfg.Queries.GetTagsByEventID(req.Context(), eventID)
	if err != nil {
		log.Printf("Error finding tags for given eventID: %s", err)
		w.WriteHeader(500)
		return
	}
	tags := []Tag{}
	for _, t := range dbTags {
		tags = append(tags, Tag{ID: t.ID, Name: t.Name, Color: t.Color})
	}
	respondWithJSON(w, 200, tags)
}

func (cfg *ApiConfig) UpdateTag(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	tagID, err := uuid.Parse(req.PathValue("tag_id"))
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

	dbTagParams := database.UpdateTagParams{Name: params.Name, Color: params.Color, TagID: tagID}
	dbTag, err := cfg.Queries.UpdateTag(req.Context(), dbTagParams)
	if err != nil {
		log.Printf("Error updating tag: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tag := Tag{ID: dbTag.ID, UserID: dbTag.UserID, Name: dbTag.Name, Color: dbTag.Color}
	respondWithJSON(w, 201, tag)
}

func (cfg *ApiConfig) DeleteTag(w http.ResponseWriter, req *http.Request) {
	tagID, err := uuid.Parse(req.PathValue("tag_id"))
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

	err = cfg.Queries.DeleteTag(req.Context(), tagID)
	if err != nil {
		log.Printf("Error deleting tag: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}

func (cfg *ApiConfig) DeleteEventTag(w http.ResponseWriter, req *http.Request) {
	tagID, err := uuid.Parse(req.PathValue("tag_id"))
	if err != nil {
		log.Printf("Error parsing uuid: %s", err)
		w.WriteHeader(500)
		return
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

	_, err = auth.ValidateAccessToken(accessToken, cfg.Secret)
	if err != nil {
		log.Printf("Access token invalid: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	deleteParams := database.DeleteEventTagParams{EventID: eventID, TagID: tagID}
	err = cfg.Queries.DeleteEventTag(req.Context(), deleteParams)
	if err != nil {
		log.Printf("Error deleting event tag: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}
