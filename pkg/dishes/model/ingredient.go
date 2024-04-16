package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Ingredient struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	DishID    string `json:"dishId"`
}

type IngredientModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (i IngredientModel) Insert(ingredient *Ingredient) error {
	fmt.Println(ingredient.Name, ingredient.Quantity)

	query := `
		INSERT INTO ingredients (name, quantity, dish_id)
		VALUES ($1, $2, $3)
		RETURNING id, createdat, updatedat
	`
	args := []interface{}{ingredient.Name, ingredient.Quantity, ingredient.DishID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return i.DB.QueryRowContext(ctx, query, args...).Scan(&ingredient.ID, &ingredient.CreatedAt, &ingredient.UpdatedAt)
}

func (i IngredientModel) GetById(id string) (*Ingredient, error) {
	query := `
		SELECT id, createdat, updatedat, name, quantity, dish_id
		FROM ingredients
		WHERE id = $1
	`
	var ingredient Ingredient
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := i.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&ingredient.ID, &ingredient.CreatedAt, &ingredient.UpdatedAt, &ingredient.Name, &ingredient.Quantity, &ingredient.DishID)

	if err != nil {
		return nil, err
	}

	return &ingredient, nil
}

func (i IngredientModel) Update(ingredient *Ingredient) error {
	query := `
		UPDATE ingredients
		SET name = $1, quantity = $2, dish_id = $3
		WHERE id = $4
		RETURNING updatedat
	`

	args := []interface{}{ingredient.Name, ingredient.Quantity, ingredient.DishID, ingredient.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return i.DB.QueryRowContext(ctx, query, args...).Scan(&ingredient.UpdatedAt)
}

func (i IngredientModel) Delete(id string) error {
	query := `
		DELETE FROM ingredients
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := i.DB.ExecContext(ctx, query, id)

	return err
}
