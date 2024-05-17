package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/validator"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	category := &model.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	err = app.models.Categories.Insert(category)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	app.respondWithJson(w, http.StatusCreated, category)
}

func (app *application) getAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		model.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{
		"id", "name",
		"-id", "-name",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	categories, metadata, err := app.models.Categories.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCategoryByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["categoryId"]

	category, err := app.models.Categories.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Category not found")
		return
	}
	app.respondWithJson(w, http.StatusOK, category)
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["categoryId"]

	category, err := app.models.Categories.GetById(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.Description != nil {
		category.Description = *input.Description
	}

	err = app.models.Categories.Update(category)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	app.respondWithJson(w, http.StatusOK, category)
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["categoryId"]

	err := app.models.Categories.Delete(param)
	if err != nil {
		app.respondWithError(w,http.StatusInternalServerError, "Failed to delete category")
		return
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
