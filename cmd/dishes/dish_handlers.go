package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJson(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createDishHandler(w http.ResponseWriter, r *http.Request) {
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

	dish := &model.Dish{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	}

	err = app.models.Dishes.Insert(dish)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to create dish")
		return
	}

	app.respondWithJson(w, http.StatusCreated, dish)
}

func (app *application) getAllDishesHandler(w http.ResponseWriter, r *http.Request) {
	dishes, err := app.models.Dishes.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch dishes")
		return
	}
	app.respondWithJson(w, http.StatusOK, dishes)
}

func (app *application) getDishByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["dishId"]

	dish, err := app.models.Dishes.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Dish not found")
		return
	}
	app.respondWithJson(w, http.StatusOK, dish)
}

func (app *application) updateDishHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["dishId"]

	dish, err := app.models.Dishes.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Dish not found")
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
		dish.Name = *input.Name
	}
	if input.Description != nil {
		dish.Description = *input.Description
	}
	if input.Price != nil {
		dish.Price = *input.Price
	}

	err = app.models.Dishes.Update(dish)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to update dish")
		return
	}

	app.respondWithJson(w, http.StatusOK, dish)
}

func (app *application) deleteDishHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["dishId"]

	err := app.models.Dishes.Delete(param)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to delete dish")
		return
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"message": "Dish deleted successfully"})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
