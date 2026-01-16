// Package handler contains HTTP request handlers for the global App.
package handler

import (
	_ "embed"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/tmm6907/mashboard/db"
)

type AppHandler struct {
	db      *sqlx.DB
	dbMutex sync.RWMutex
}

func NewAppHandler() AppHandler {
	db := db.InitDB()
	return AppHandler{db, sync.RWMutex{}}
}
