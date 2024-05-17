package model

import (
	"database/sql"
	"log"
	"time"
)

// Review struct represents a review.
type Review struct {
	ID        int       `json:"id"`
	DishID    *int      `json:"dish_id,omitempty"`
	DrinkID   *int      `json:"drink_id,omitempty"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ReviewInput struct represents the input data for creating or updating a review.
type ReviewInput struct {
	DishID  *int   `json:"dish_id,omitempty"`
	DrinkID *int   `json:"drink_id,omitempty"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// ReviewModel struct wraps a sql.DB connection pool.
type ReviewModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// Create inserts a new review into the database.
func (m ReviewModel) Create(input ReviewInput) (*Review, error) {
	query := `
		INSERT INTO reviews (dish_id, drink_id, rating, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, dish_id, drink_id, rating, comment, created_at, updated_at`

	row := m.DB.QueryRow(query, input.DishID, input.DrinkID, input.Rating, input.Comment)

	var review Review
	err := row.Scan(&review.ID, &review.DishID, &review.DrinkID, &review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

// GetByID fetches a review by its ID.
func (m ReviewModel) GetByID(id int) (*Review, error) {
	query := `
		SELECT id, dish_id, drink_id, rating, comment, created_at, updated_at
		FROM reviews
		WHERE id = $1`

	var review Review
	err := m.DB.QueryRow(query, id).Scan(&review.ID, &review.DishID, &review.DrinkID, &review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &review, nil
}

// Update modifies an existing review in the database.
func (m ReviewModel) Update(id int, input ReviewInput) (*Review, error) {
	query := `
		UPDATE reviews
		SET dish_id = $1, drink_id = $2, rating = $3, comment = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, dish_id, drink_id, rating, comment, created_at, updated_at`

	row := m.DB.QueryRow(query, input.DishID, input.DrinkID, input.Rating, input.Comment, id)

	var review Review
	err := row.Scan(&review.ID, &review.DishID, &review.DrinkID, &review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &review, nil
}

// Delete removes a review from the database by its ID.
func (m ReviewModel) Delete(id int) error {
	query := `DELETE FROM reviews WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// GetAll fetches all reviews with optional filtering by entity ID and supports pagination.
func (m ReviewModel) GetAll(dishID, drinkID *int, filters Filters) ([]*Review, Metadata, error) {
	query := `
		SELECT COUNT(*) OVER(), id, dish_id, drink_id, rating, comment, created_at, updated_at
		FROM reviews
		WHERE (dish_id = $1 OR $1 IS NULL) AND (drink_id = $2 OR $2 IS NULL)
		ORDER BY ` + filters.sortColumn() + ` ` + filters.sortDirection() + `
		LIMIT $3 OFFSET $4`

	args := []interface{}{dishID, drinkID, filters.limit(), filters.offset()}

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	reviews := []*Review{}
	totalRecords := 0

	for rows.Next() {
		var review Review
		err := rows.Scan(
			&totalRecords,
			&review.ID,
			&review.DishID,
			&review.DrinkID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		reviews = append(reviews, &review)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return reviews, metadata, nil
}
