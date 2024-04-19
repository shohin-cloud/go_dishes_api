package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/validator"
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
	dishes, metadata, err := app.models.Dishes.GetAll(input.Name, input.Price, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"dishes": dishes, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope,
	headers http.Header) error {
	// Use the json.MarshalIndent() function so that whitespace is added to the encoded JSON. Use
	// no line prefix and tab indents for each element.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// At this point, we know that we won't encounter any more errors before writing the response,
	// so it's safe to add any headers that we want to include. We loop through the header map
	// and add each header to the http.ResponseWriter header map. Note that it's OK if the
	// provided header map is nil. Go doesn't through an error if you try to range over (
	// or generally, read from) a nil map
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(js); err != nil {
		app.logger.PrintError(err, nil)
		return err
	}

	return nil
}
