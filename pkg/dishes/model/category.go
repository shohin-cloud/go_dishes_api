package model

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"
)

type Category struct {
    ID          string    `json:"id"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
}

type CategoryModel struct {
    DB       *sql.DB
    InfoLog  *log.Logger
    ErrorLog *log.Logger
}

func (c CategoryModel) Insert(category *Category) error {
    query := `
        INSERT INTO categories (name, description)
        VALUES ($1, $2)
        RETURNING id, createdat, updatedat
    `
    args := []interface{}{category.Name, category.Description}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.DB.QueryRowContext(ctx, query, args...).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

func (c CategoryModel) GetAll(name string, filters Filters) ([]*Category, Metadata, error) {
    query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, createdat, updatedat, name, description
        FROM categories
        WHERE (LOWER(name) LIKE LOWER($1) OR $1 = '')
        ORDER BY %s %s, id ASC
        LIMIT $2 OFFSET $3
    `, filters.sortColumn(), filters.sortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    args := []interface{}{"%" + name + "%", filters.limit(), filters.offset()}

    rows, err := c.DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    categories := []*Category{}

    for rows.Next() {
        var category Category

        err := rows.Scan(
            &totalRecords,
            &category.ID,
            &category.CreatedAt,
            &category.UpdatedAt,
            &category.Name,
            &category.Description,
        )

        if err != nil {
            return nil, Metadata{}, err
        }

        categories = append(categories, &category)
    }

    if err := rows.Err(); err != nil {
        return nil, Metadata{}, err
    }

    metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

    return categories, metadata, nil
}

func (c CategoryModel) GetById(id string) (*Category, error) {
    query := `
        SELECT id, createdat, updatedat, name, description
        FROM categories
        WHERE id = $1
    `
    var category Category
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    row := c.DB.QueryRowContext(ctx, query, id)
    err := row.Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt, &category.Name, &category.Description)

    if err != nil {
        return nil, err
    }

    return &category, nil
}

func (c CategoryModel) Update(category *Category) error {
    query := `
        UPDATE categories
        SET name = $1, description = $2
        WHERE id = $3
        RETURNING updatedat
    `
    args := []interface{}{category.Name, category.Description, category.ID}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return c.DB.QueryRowContext(ctx, query, args...).Scan(&category.UpdatedAt)
}

func (c CategoryModel) Delete(id string) error {
    query := `
        DELETE FROM categories
        WHERE id = $1
    `
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := c.DB.ExecContext(ctx, query, id)

    return err
}
