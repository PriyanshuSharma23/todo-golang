package data

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("requested resource not found")
)

type Models struct {
	TodosModel *TodosModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		TodosModel: &TodosModel{db: db},
	}
}
