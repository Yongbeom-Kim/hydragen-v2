package handler

import "database/sql"

type App struct {
	Db *sql.DB
}

func (a *App) UseFallbackData() bool {
	return a.Db == nil
}
