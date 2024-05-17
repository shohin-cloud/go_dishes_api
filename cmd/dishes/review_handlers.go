package main

import (
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/shohin-cloud/dishes-api/pkg/dishes/model"
    "github.com/shohin-cloud/dishes-api/pkg/dishes/validator"
)

// CreateReviewHandler creates a new review.
func (app *application) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
    var input model.ReviewInput
    err := app.readJSON(w, r, &input)
    if err != nil {
        app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    review, err := app.models.Reviews.Create(input)
    if err != nil {
        app.respondWithError(w, http.StatusInternalServerError, "Failed to create review")
        return
    }

    app.respondWithJson(w, http.StatusCreated, review)
}

// GetReviewHandler retrieves a review by ID.
func (app *application) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    reviewID, err := strconv.Atoi(vars["reviewID"])
    if err != nil {
        app.respondWithError(w, http.StatusBadRequest, "Invalid review ID")
        return
    }

    review, err := app.models.Reviews.GetByID(reviewID)
    if err != nil {
        app.respondWithError(w, http.StatusNotFound, "Review not found")
        return
    }

    app.respondWithJson(w, http.StatusOK, review)
}

// UpdateReviewHandler updates an existing review.
func (app *application) UpdateReviewHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    reviewID, err := strconv.Atoi(vars["reviewID"])
    if err != nil {
        app.respondWithError(w, http.StatusBadRequest, "Invalid review ID")
        return
    }

    var input model.ReviewInput
    err = app.readJSON(w, r, &input)
    if err != nil {
        app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    updatedReview, err := app.models.Reviews.Update(reviewID, input)
    if err != nil {
        app.respondWithError(w, http.StatusInternalServerError, "Failed to update review")
        return
    }

    app.respondWithJson(w, http.StatusOK, updatedReview)
}

// DeleteReviewHandler deletes a review by ID.
func (app *application) DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    reviewID, err := strconv.Atoi(vars["reviewID"])
    if err != nil {
        app.respondWithError(w, http.StatusBadRequest, "Invalid review ID")
        return
    }

    err = app.models.Reviews.Delete(reviewID)
    if err != nil {
        app.respondWithError(w, http.StatusInternalServerError, "Failed to delete review")
        return
    }

    app.respondWithJson(w, http.StatusOK, map[string]string{"message": "Review deleted successfully"})
}

// GetAllReviewsHandler retrieves all reviews with pagination, filtering, and sorting.
func (app *application) GetAllReviewsHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        DishID  *int
        DrinkID *int
        model.Filters
    }

    v := validator.New()

    qs := r.URL.Query()

    if dishIDStr := app.readString(qs, "dish_id", ""); dishIDStr != "" {
        dishID, err := strconv.Atoi(dishIDStr)
        if err == nil {
            input.DishID = &dishID
        }
    }

    if drinkIDStr := app.readString(qs, "drink_id", ""); drinkIDStr != "" {
        drinkID, err := strconv.Atoi(drinkIDStr)
        if err == nil {
            input.DrinkID = &drinkID
        }
    }

    input.Filters.Page = app.readInt(qs, "page", 1, v)
    input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
    input.Filters.Sort = app.readString(qs, "sort", "id")

    input.Filters.SortSafelist = []string{
        "id", "rating", "created_at",
        "-id", "-rating", "-created_at",
    }

    if model.ValidateFilters(v, input.Filters); !v.Valid() {
        app.failedValidationResponse(w, r, v.Errors)
        return
    }

    reviews, metadata, err := app.models.Reviews.GetAll(input.DishID, input.DrinkID, input.Filters)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"reviews": reviews, "metadata": metadata}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
