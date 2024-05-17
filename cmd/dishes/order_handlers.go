package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/validator"
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID     string                `json:"userId"`
		TotalPrice float64               `json:"totalPrice"`
		Status     string                `json:"status"`
		OrderItems []model.OrderItem `json:"orderItems"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	order := &model.Order{
		UserID:     input.UserID,
		TotalPrice: input.TotalPrice,
		Status:     input.Status,
	}

	for _, item := range input.OrderItems {
		order.OrderItems = append(order.OrderItems, model.OrderItem{
			DishID:   item.DishID,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}

	err = app.models.Orders.Insert(order)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	app.respondWithJson(w, http.StatusCreated, order)
}

func (app *application) getAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Status string
		model.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Status = app.readString(qs, "status", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{
		"id", "user_id", "total_price", "status",
		"-id", "-user_id", "-total_price", "-status",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	orders, metadata, err := app.models.Orders.GetAll(input.Status, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getOrderByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

	order, err := app.models.Orders.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}
	app.respondWithJson(w, http.StatusOK, order)
}

func (app *application) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

	order, err := app.models.Orders.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	var input struct {
		UserID     *string                `json:"userId"`
		TotalPrice *float64               `json:"totalPrice"`
		Status     *string                `json:"status"`
		OrderItems []model.OrderItemInput `json:"orderItems"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.UserID != nil {
		order.UserID = *input.UserID
	}
	if input.TotalPrice != nil {
		order.TotalPrice = *input.TotalPrice
	}
	if input.Status != nil {
		order.Status = *input.Status
	}

	if input.OrderItems != nil {
		order.OrderItems = nil
		for _, item := range input.OrderItems {
			order.OrderItems = append(order.OrderItems, model.OrderItem{
				OrderID: order.ID,
				DishID:  item.DishID,
				Quantity: item.Quantity,
				Price:    item.Price,
			})
		}
	}

	err = app.models.Orders.Update(order)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to update order")
		return
	}

	app.respondWithJson(w, http.StatusOK, order)
}

func (app *application) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

	err := app.models.Orders.Delete(param)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to delete order")
		return
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"message": "Order deleted successfully"})
}
