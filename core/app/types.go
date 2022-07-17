package app

import "database/sql"

type AppRepository struct {
	db *sql.DB
}

type StartAppResponse struct {
	Success bool
	Address string
}
