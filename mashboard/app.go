package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "embed"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed build.sql
var buildScript string

// App struct
type App struct {
	ctx     context.Context
	db      *sqlx.DB
	dbMutex sync.RWMutex
}

func initDB() *sqlx.DB {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	dbDir := filepath.Join(baseDir, "mashboard")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		panic(err)
	}

	db := sqlx.MustOpen("sqlite3", filepath.Join(dbDir, "root.db"))
	db.MustExec(buildScript)

	return db
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		db: initDB(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	interval := time.Minute * 5
	a.StartRSSFetcher(&interval)
}
