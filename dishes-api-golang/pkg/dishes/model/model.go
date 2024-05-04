package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

type Models struct {
	Dishes      DishModel
	Ingredients IngredientModel
	Members     MemberModel
	Tokens      TokenModel
	Permissions PermissionModel
}

var (
	// ErrRecordNotFound is returned when a record doesn't exist in database.
	ErrRecordNotFound = errors.New("record not found")

	// ErrEditConflict is returned when a there is a data race, and we have an edit conflict.
	ErrEditConflict = errors.New("edit conflict")
)

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Dishes: DishModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Ingredients: IngredientModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Members: MemberModel{
			DB: db,
		},
		Permissions: PermissionModel{
			DB: db,
		},
		Tokens: TokenModel{
			DB: db,
		},
	}
}
