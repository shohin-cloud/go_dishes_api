package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Ingredient struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	DishID   string `json:"dishId"`
}

type IngredientModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (i IngredientModel) Insert(ingredient *Ingredient) error {
	query := `
		INSERT INTO ingredients (name, quantity, dish_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	args := []interface{}{ingredient.Name, ingredient.Quantity, ingredient.DishID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return i.DB.QueryRowContext(ctx, query, args...).Scan(&ingredient.ID)
}

func (i IngredientModel) GetAllByDishID(dishID string) ([]*Ingredient, error) {
	query := `
		SELECT id, name, quantity, dish_id
		FROM ingredients
		WHERE dish_id = $1
		ORDER BY id
	`

	rows, err := i.DB.Query(query, dishID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []*Ingredient
	for rows.Next() {
		var ingredient Ingredient
		err := rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Quantity, &ingredient.DishID)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, &ingredient)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ingredients, nil
}

func (i IngredientModel) DeleteByDishID(dishID string) error {
	query := `
		DELETE FROM ingredients
		WHERE dish_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := i.DB.ExecContext(ctx, query, dishID)

	return err
}
