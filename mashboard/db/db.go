package db

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

//go:embed build.sql
var buildScript string

func InitDB() *sqlx.DB {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	dbDir := filepath.Join(baseDir, "mashboard")
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		panic(err)
	}

	db := sqlx.MustOpen("sqlite3", filepath.Join(dbDir, "root.db"))
	db.MustExec(buildScript)

	return db
}
