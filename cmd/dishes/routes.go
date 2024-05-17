package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// routes is our main application's router.
func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck
	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Dishes
	// v1.HandleFunc("/dishes", app.requirePermission("dishes:write", app.createDishHandler)).Methods("POST")
	// v1.HandleFunc("/dishes", app.requirePermission("dishes:read", app.getAllDishesHandler)).Methods("GET")
	// v1.HandleFunc("/dishes/{dishId:[0-9]+}", app.requirePermission("dishes:read", app.getDishByIdHandler)).Methods("GET")
	// v1.HandleFunc("/dishes/{dishId:[0-9]+}", app.requirePermission("dishes:write", app.updateDishHandler)).Methods("PUT")
	// v1.HandleFunc("/dishes/{dishId:[0-9]+}", app.requirePermission("dishes:write", app.deleteDishHandler)).Methods("DELETE")

	v1.HandleFunc("/dishes", app.createDishHandler).Methods("POST")
	v1.HandleFunc("/dishes", app.getAllDishesHandler).Methods("GET")
	v1.HandleFunc("/dishes/{dishId:[0-9]+}", app.getDishByIdHandler).Methods("GET")
	v1.HandleFunc("/dishes/{dishId:[0-9]+}", app.updateDishHandler).Methods("PUT")
	v1.HandleFunc("/dishes/{dishId:[0-9]+}", app.deleteDishHandler).Methods("DELETE")

    // Drinks
    v1.HandleFunc("/drinks", app.createDrinkHandler).Methods("POST")
    v1.HandleFunc("/drinks", app.getAllDrinksHandler).Methods("GET")
    v1.HandleFunc("/drinks/{drinkId:[0-9]+}", app.getDrinkByIdHandler).Methods("GET")
    v1.HandleFunc("/drinks/{drinkId:[0-9]+}", app.updateDrinkHandler).Methods("PUT")
    v1.HandleFunc("/drinks/{drinkId:[0-9]+}", app.deleteDrinkHandler).Methods("DELETE")

	// Ingredients
	// v1.HandleFunc("/ingredients", app.createIngredientHandler).Methods("POST")
	// v1.HandleFunc("/ingredients/{ingredientId:[0-9]+}", app.getIngredientByIdHandler).Methods("GET")
	// v1.HandleFunc("/ingredients/{ingredientId:[0-9]+}", app.updateIngredientHandler).Methods("PUT")
	// v1.HandleFunc("/ingredients/{ingredientId:[0-9]+}", app.deleteIngredientHandler).Methods("DELETE")

	v1.HandleFunc("/ingredients", app.createIngredientHandler).Methods("POST")
	v1.HandleFunc("/ingredients/{ingredientId:[0-9]+}", app.getIngredientByIdHandler).Methods("GET")
	v1.HandleFunc("/ingredients/{ingredientId:[0-9]+}", app.updateIngredientHandler).Methods("PUT")
	v1.HandleFunc("/ingredients/{ingredientId:[0-9]+}", app.deleteIngredientHandler).Methods("DELETE")
	// Members
	v1.HandleFunc("/members", app.registerMemberHandler).Methods("POST")
	v1.HandleFunc("/members/activated", app.activateMemberHandler).Methods("PUT")
	v1.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
