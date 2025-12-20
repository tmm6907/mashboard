package main

func (a *App) CreateNewCollection(name string) Response {
	if name == "" {
		return Response{"name must not be empty", nil}
	}
	var i Collection
	if err := a.db.Get(&i, "SELECT * FROM collections WHERE name = ?", name); err == nil {
		return Response{"collection already exists", nil}
	}
	if _, err := a.db.Exec("INSERT INTO collections (name) VALUES (?);", name); err != nil {
		return Response{err.Error(), nil}
	}
	return Response{"", "Success"}
}

func (a *App) GetCollections() Response {
	var collections []Collection
	if err := a.db.Select(&collections, "SELECT * FROM collections;"); err != nil {
		return Response{err.Error(), nil}
	}
	return Response{"", collections}
}
