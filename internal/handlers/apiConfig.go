package handlers

import "github.com/curtisbraxdale/taday/internal/database"

type ApiConfig struct {
	Queries  *database.Queries
	Platform string
	Secret   string
}
