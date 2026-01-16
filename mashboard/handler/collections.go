package handler

import (
	"errors"

	"github.com/tmm6907/mashboard/models"
)

func (a *AppHandler) CreateNewCollection(name string) Response {
	if name == "" {
		return a.SendError(errors.New("name must not be empty"))
	}
	var i models.Collection
	if err := a.db.Get(&i, "SELECT * FROM collections WHERE name = ?", name); err == nil {
		return a.SendError(errors.New("collection already exists"))
	}
	if _, err := a.db.Exec("INSERT INTO collections (name) VALUES (?);", name); err != nil {
		return a.SendError(err)
	}
	return a.SendResponse("Success")
}

func (a *AppHandler) GetCollections() Response {
	var collections []models.Collection
	if err := a.db.Select(&collections, "SELECT * FROM collections;"); err != nil {
		return a.SendError(err)
	}
	return a.SendResponse(collections)
}
