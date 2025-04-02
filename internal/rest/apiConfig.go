package rest

import "github.com/hconn7/GoGive/internal/database"

type ApiConfig struct {
	DbQueries   *database.Queries
	TokenSecret string
}
