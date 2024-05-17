package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Drink represents a drink entity.
type Drink struct {
	ID          string  `json:"id"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// DrinkModel manages interactions with the drink table in the database.
type DrinkModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// Insert inserts a new drink into the database.
func (d DrinkModel) Insert(drink *Drink) error {
	fmt.Println(drink.Name, drink.Description, drink.Price)

	query := `
		INSERT INTO drinks (name, description, price)
		VALUES ($1, $2, $3)
		RETURNING id, createdat, updatedat
	`
	args := []interface{}{drink.Name, drink.Description, drink.Price}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return d.DB.QueryRowContext(ctx, query, args...).Scan(&drink.ID, &drink.CreatedAt, &drink.UpdatedAt)
}

// GetAll retrieves all drinks from the database.
func (d DrinkModel) GetAll(name string, price string, filters Filters) ([]*Drink, Metadata, error) {
	if price == "" {
		price = "0"
	}
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, createdAt, updatedAt, name, description, price
		FROM drinks
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
	drinks := []*Drink{}

	for rows.Next() {
		var drink Drink

		err := rows.Scan(
			&totalRecords,
			&drink.ID,
			&drink.CreatedAt,
			&drink.UpdatedAt,
			&drink.Name,
			&drink.Description,
			&drink.Price,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		drinks = append(drinks, &drink)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return drinks, metadata, nil
}

// GetById retrieves a drink by ID from the database.
func (d DrinkModel) GetById(id string) (*Drink, error) {
	query := `
		SELECT id, createdat, updatedat, name, description, price
		FROM drinks
		WHERE id = $1
	`
	var drink Drink
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := d.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&drink.ID, &drink.CreatedAt, &drink.UpdatedAt, &drink.Name, &drink.Description, &drink.Price)

	if err != nil {
		return nil, err
	}

	return &drink, nil
}

// Update updates a drink in the database.
func (d DrinkModel) Update(drink *Drink) error {
	query := `
		UPDATE drinks
		SET name = $1, description = $2, price = $3
		WHERE id = $4
		RETURNING updatedat
	`

	args := []interface{}{drink.Name, drink.Description, drink.Price, drink.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return d.DB.QueryRowContext(ctx, query, args...).Scan(&drink.UpdatedAt)
}

// Delete deletes a drink from the database.
func (d DrinkModel) Delete(id string) error {
	query := `
		DELETE FROM drinks
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := d.DB.ExecContext(ctx, query, id)

	return err
}
