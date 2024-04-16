package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shohin-cloud/dishes-api/pkg/dishes/model"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

func main() {
	fmt.Println("Started server")
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:admin@localhost/go_project?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func (app *application) run() {
	fmt.Println("Running")
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Dishes
	v1.HandleFunc("/dishes", app.createDishHandler).Methods("POST")
	v1.HandleFunc("/dishes", app.getAllDishesHandler).Methods("GET")
	v1.HandleFunc("/dishes/{dishId:[0-9]+}", app.getDishByIdHandler).Methods("GET")
	v1.HandleFunc("/dishes/{dishId:[0-9]+}", app.updateDishHandler).Methods("PUT")
	v1.HandleFunc("/dishes/{dishId:[0-9]+}", app.deleteDishHandler).Methods("DELETE")

	// Ingredients
	v1.HandleFunc("/ingredients", app.createIngredientHandler).Methods("POST")
	v1.HandleFunc("/ingredients/{ingredientId:[0-9]+}", app.getIngredientByIdHandler).Methods("GET")
	v1.HandleFunc("/ingredients/{ingredientId:[0-9]+}", app.updateIngredientHandler).Methods("PUT")
	v1.HandleFunc("/ingredients/{ingredientId:[0-9]+}", app.deleteIngredientHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
