package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) CreateDailyAgenda(userID uuid.UUID) (string, error) {
	now := time.Now()
	dbToDos, err := cfg.Queries.GetTodosByUserID(context.Background(), userID)
	if err != nil {
		log.Printf("Error getting ToDos for user %v: %s", userID, err)
		return "", err
	}
	dbEventParams := database.GetUserEventsTodayParams{UserID: userID, Date: now, DatePlus1Day: now.Add(24 * time.Hour)}
	dbEvents, err := cfg.Queries.GetUserEventsToday(context.Background(), dbEventParams)
	if err != nil {
		log.Printf("Error getting Events for user %v: %s", userID, err)
		return "", err
	}
	agendaString := "=== TADAYs AGENDA ===\n\n"
	if len(dbToDos) > 0 {
		for _, t := range dbToDos {
			if t.Description.Valid {
				agendaString += fmt.Sprintf("%v\n+%v\n\n", t.Title, t.Description.String)
			} else {
				agendaString += fmt.Sprintf("%v\n\n", t.Title)
			}
		}
	}
	if len(dbEvents) > 0 {
		for _, e := range dbEvents {
			if e.Description.Valid {
				agendaString += fmt.Sprintf("%v\n+%v\n\n", e.Title, e.Description.String)
			} else {
				agendaString += fmt.Sprintf("%v\n\n", e.Title)
			}
		}
	}
	agendaString += fmt.Sprint("=====================\n")
	return agendaString, nil
}

func (cfg *ApiConfig) CreateWeeklyAgenda(userID uuid.UUID) (string, error) {
	now := time.Now()
	dbToDos, err := cfg.Queries.GetTodosByUserID(context.Background(), userID)
	if err != nil {
		log.Fatalf("Error getting ToDos for user %v: %s", userID, err)
		return "", err
	}
	dbEventParams := database.GetUserEventsWeekParams{UserID: userID, Date: now, DatePlus7Days: now.Add(7 * 24 * time.Hour)}
	dbEvents, err := cfg.Queries.GetUserEventsWeek(context.Background(), dbEventParams)
	if err != nil {
		log.Fatalf("Error getting Events for user %v: %s", userID, err)
		return "", err
	}
	agendaString := "=== TADAYs AGENDA ===\n\n"
	for _, t := range dbToDos {
		if t.Description.Valid {
			agendaString += fmt.Sprintf("%v\n+%v\n\n", t.Title, t.Description.String)
		} else {
			agendaString += fmt.Sprintf("%v\n\n", t.Title)
		}
	}
	for _, e := range dbEvents {
		if e.Description.Valid {
			agendaString += fmt.Sprintf("%v\n+%v\n\n", e.Title, e.Description.String)
		} else {
			agendaString += fmt.Sprintf("%v\n\n", e.Title)
		}
	}
	agendaString += fmt.Sprint("=====================\n")
	return agendaString, nil
}
