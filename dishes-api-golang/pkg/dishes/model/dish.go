package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Dish struct {
	ID          string  `json:"id"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type DishModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (d DishModel) Insert(dish *Dish) error {
	fmt.Println(dish.Name, dish.Description, dish.Price)

	query := `
		INSERT INTO dishes (name, description, price)
		VALUES ($1, $2, $3)
		RETURNING id, createdat, updatedat
	`
	args := []interface{}{dish.Name, dish.Description, dish.Price}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return d.DB.QueryRowContext(ctx, query, args...).Scan(&dish.ID, &dish.CreatedAt, &dish.UpdatedAt)
}

func (d DishModel) GetAll(name string, price string, filters Filters) ([]*Dish, Metadata, error) {
	if price == "" {
		price = "0"
	}
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, createdAt, updatedAt, name, description, price
		FROM dishes
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (price >= $2 OR $2 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, price, filters.limit(), filters.offset()}

	rows, err := d.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	dishes := []*Dish{}

	for rows.Next() {
		var dish Dish

		err := rows.Scan(
			&totalRecords,
			&dish.ID,
			&dish.CreatedAt,
			&dish.UpdatedAt,
			&dish.Name,
			&dish.Description,
			&dish.Price,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		dishes = append(dishes, &dish)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return dishes, metadata, nil // Return users and nil error
}

func (d DishModel) GetById(id string) (*Dish, error) {
	query := `
		SELECT id, createdat, updatedat, name, description, price
		FROM dishes
		WHERE id = $1
	`
	var dish Dish
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := d.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&dish.ID, &dish.CreatedAt, &dish.UpdatedAt, &dish.Name, &dish.Description, &dish.Price)

	if err != nil {
		return nil, err
	}

	return &dish, nil
}

func (d DishModel) Update(dish *Dish) error {
	query := `
		UPDATE dishes
		SET name = $1, description = $2, price = $3
		WHERE id = $4
		RETURNING updatedat
	`

	args := []interface{}{dish.Name, dish.Description, dish.Price, dish.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return d.DB.QueryRowContext(ctx, query, args...).Scan(&dish.UpdatedAt)
}

func (d DishModel) Delete(id string) error {
	query := `
		DELETE FROM dishes
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := d.DB.ExecContext(ctx, query, id)

	return err
}
