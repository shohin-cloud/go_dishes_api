package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
)

func (app *application) createIngredientHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		DishID   string `json:"dishId"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ingredient := &model.Ingredient{
		Name:     input.Name,
		Quantity: input.Quantity,
		DishID:   input.DishID,
	}

	err = app.models.Ingredients.Insert(ingredient)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to create ingredient")
		return
	}

	app.respondWithJson(w, http.StatusCreated, ingredient)
}

func (app *application) getIngredientByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["ingredientId"]

	ingredient, err := app.models.Ingredients.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Ingredient not found")
		return
	}
	app.respondWithJson(w, http.StatusOK, ingredient)
}

func (app *application) updateIngredientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["ingredientId"]

	ingredient, err := app.models.Ingredients.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Ingredient not found")
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Quantity *int    `json:"quantity"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		ingredient.Name = *input.Name
	}
	if input.Quantity != nil {
		ingredient.Quantity = *input.Quantity
	}

	err = app.models.Ingredients.Update(ingredient)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to update ingredient")
		return
	}

	app.respondWithJson(w, http.StatusOK, ingredient)
}

func (app *application) deleteIngredientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["ingredientId"]

	err := app.models.Ingredients.Delete(param)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to delete ingredient")
		return
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"message": "Ingredient deleted successfully"})
}
