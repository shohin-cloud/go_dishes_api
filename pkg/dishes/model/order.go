package model

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"
)

type Order struct {
    ID         string       `json:"id"`
    CreatedAt  time.Time    `json:"createdAt"`
    UpdatedAt  time.Time    `json:"updatedAt"`
    UserID     string       `json:"userId"`
    TotalPrice float64      `json:"totalPrice"`
    Status     string       `json:"status"`
    OrderItems []OrderItem  `json:"orderItems"`
}

type OrderItem struct {
    ID        string  `json:"id"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    OrderID   string  `json:"orderId"`
    DishID    string  `json:"dishId"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
}

type OrderModel struct {
    DB       *sql.DB
    InfoLog  *log.Logger
    ErrorLog *log.Logger
}

func (o OrderModel) Insert(order *Order) error {
    query := `
        INSERT INTO orders (user_id, total_price, status)
        VALUES ($1, $2, $3)
        RETURNING id, createdat, updatedat
    `
    args := []interface{}{order.UserID, order.TotalPrice, order.Status}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := o.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
    if err != nil {
        return err
    }

    for _, item := range order.OrderItems {
        item.OrderID = order.ID
        err := o.InsertOrderItem(&item)
        if err != nil {
            return err
        }
    }

    return nil
}

func (o OrderModel) InsertOrderItem(orderItem *OrderItem) error {
    query := `
        INSERT INTO order_items (order_id, dish_id, quantity, price)
        VALUES ($1, $2, $3, $4)
        RETURNING id, createdat, updatedat
    `
    args := []interface{}{orderItem.OrderID, orderItem.DishID, orderItem.Quantity, orderItem.Price}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return o.DB.QueryRowContext(ctx, query, args...).Scan(&orderItem.ID, &orderItem.CreatedAt, &orderItem.UpdatedAt)
}

func (o OrderModel) GetAll(status string, filters Filters) ([]*Order, Metadata, error) {
    query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, createdat, updatedat, user_id, total_price, status
        FROM orders
        WHERE (LOWER(status) = LOWER($1) OR $1 = '')
        ORDER BY %s %s, id ASC
        LIMIT $2 OFFSET $3
    `, filters.sortColumn(), filters.sortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    args := []interface{}{status, filters.limit(), filters.offset()}

    rows, err := o.DB.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    orders := []*Order{}

    for rows.Next() {
        var order Order

        err := rows.Scan(
            &totalRecords,
            &order.ID,
            &order.CreatedAt,
            &order.UpdatedAt,
            &order.UserID,
            &order.TotalPrice,
            &order.Status,
        )

        if err != nil {
            return nil, Metadata{}, err
        }

        order.OrderItems, err = o.GetOrderItems(order.ID)
        if err != nil {
            return nil, Metadata{}, err
        }

        orders = append(orders, &order)
    }

    if err := rows.Err(); err != nil {
        return nil, Metadata{}, err
    }

    metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

    return orders, metadata, nil
}

func (o OrderModel) GetOrderItems(orderID string) ([]OrderItem, error) {
    query := `
        SELECT id, createdat, updatedat, order_id, dish_id, quantity, price
        FROM order_items
        WHERE order_id = $1
    `
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := o.DB.QueryContext(ctx, query, orderID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    orderItems := []OrderItem{}
    for rows.Next() {
        var item OrderItem
        err := rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.OrderID, &item.DishID, &item.Quantity, &item.Price)
        if err != nil {
            return nil, err
        }
        orderItems = append(orderItems, item)
    }
    return orderItems, rows.Err()
}

func (o OrderModel) GetById(id string) (*Order, error) {
    query := `
        SELECT id, createdat, updatedat, user_id, total_price, status
        FROM orders
        WHERE id = $1
    `
    var order Order
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    row := o.DB.QueryRowContext(ctx, query, id)
    err := row.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.UserID, &order.TotalPrice, &order.Status)
    if err != nil {
        return nil, err
    }

    order.OrderItems, err = o.GetOrderItems(order.ID)
    if err != nil {
        return nil, err
    }

    return &order, nil
}

func (o OrderModel) Update(order *Order) error {
    query := `
        UPDATE orders
        SET user_id = $1, total_price = $2, status = $3
        WHERE id = $4
        RETURNING updatedat
    `
    args := []interface{}{order.UserID, order.TotalPrice, order.Status, order.ID}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := o.DB.QueryRowContext(ctx, query, args...).Scan(&order.UpdatedAt)
    if err != nil {
        return err
    }

    // Update order items if needed
    for _, item := range order.OrderItems {
        err := o.UpdateOrderItem(&item)
        if err != nil {
            return err
        }
    }

    return nil
}

func (o OrderModel) UpdateOrderItem(orderItem *OrderItem) error {
    query := `
        UPDATE order_items
        SET dish_id = $1, quantity = $2, price = $3
        WHERE id = $4
        RETURNING updatedat
    `
    args := []interface{}{orderItem.DishID, orderItem.Quantity, orderItem.Price, orderItem.ID}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return o.DB.QueryRowContext(ctx, query, args...).Scan(&orderItem.UpdatedAt)
}

func (o OrderModel) Delete(id string) error {
    query := `
        DELETE FROM orders
        WHERE id = $1
    `
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := o.DB.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    return o.DeleteOrderItems(id)
}

func (o OrderModel) DeleteOrderItems(orderID string) error {
    query := `
        DELETE FROM order_items
        WHERE order_id = $1
    `
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := o.DB.ExecContext(ctx, query, orderID)

    return err
}
