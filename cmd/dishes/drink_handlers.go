package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/validator"
)

func (app *application) createDrinkHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	drink := &model.Drink{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	}

	err = app.models.Drinks.Insert(drink)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to create drink")
		return
	}

	app.respondWithJson(w, http.StatusCreated, drink)
}

func (app *application) getAllDrinksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string
		Price string
		model.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Price = app.readString(qs, "price", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{
		"id", "name", "price",
		"-id", "-name", "-price",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	drinks, metadata, err := app.models.Drinks.GetAll(input.Name, input.Price, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"drinks": drinks, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getDrinkByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["drinkId"]

	drink, err := app.models.Drinks.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Drink not found")
		return
	}
	app.respondWithJson(w, http.StatusOK, drink)
}

func (app *application) updateDrinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["drinkId"]

	drink, err := app.models.Drinks.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Drink not found")
		return
	}

	var input struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Price       *float64 `json:"price"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		drink.Name = *input.Name
	}
	if input.Description != nil {
		drink.Description = *input.Description
	}
	if input.Price != nil {
		drink.Price = *input.Price
	}

	err = app.models.Drinks.Update(drink)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to update drink")
		return
	}

	app.respondWithJson(w, http.StatusOK, drink)
}

func (app *application) deleteDrinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["drinkId"]

	err := app.models.Drinks.Delete(param)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to delete drink")
		return
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"message": "Drink deleted successfully"})
}
