# Dishes API

This project is a Go-based API server for managing dishes. It connects to a PostgreSQL database, uses a custom logging package, and supports different environments (development, staging, production).

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Server](#running-the-server)
- [Project Structure](#project-structure)

## Installation

Before you start, ensure you have the following installed:

- [Go](https://golang.org/doc/install) (version 1.16+)
- [PostgreSQL](https://www.postgresql.org/download/)

### Cloning the Repository

Clone this repository to your local machine using:

```sh
git clone https://github.com/shohin-cloud/dishes-api.git
cd dishes-api
```

## Installing Dependencies

Navigate to the project directory and run:

```sh
go mod tidy
```

## Configuration

The server configuration is managed through command-line flags. These flags configure the server's port, environment, and database connection string.

### Command-Line Flags

| Flag       | Default Value                                                                                                    | Description                                 |
|------------|------------------------------------------------------------------------------------------------------------------|---------------------------------------------|
| `-port`    | `:8060`                                                                                                          | The port on which the server will listen.   |
| `-env`     | `development`                                                                                                    | The application environment (development, staging, production). |
| `-db-dsn`  | `postgres://dishes_db_user:r7Y1IaUPdozGB8o27WmOjtysdS2aoBHN@dpg-cp28gs779t8c73fotm60-a.singapore-postgres.render.com/dishes_db` | The PostgreSQL DSN for connecting to the database. |

### Example Usage

You can start the server with the default configuration by running:

```sh
go run .
```

# Dishes REST API

```sh
# list of all menu items
GET /dishes
Retrieve a list of all dishes items.

POST /dishes
Create a new dish item.

GET /dishes/:id
Retrieve a specific dish item by its ID.

PUT /dishes/:id
Update an existing dish item by its ID.

DELETE /dishes/:id
Delete a specific dish item by its ID.
```
# Ingredients REST API

```sh
# list of all menu items
GET /ingredients
Retrieve a list of all ingredients items.

POST /ingredients
Create a new ingredients item.

GET /ingredients/:id
Retrieve a specific ingredients item by its ID.

PUT /ingredients/:id
Update an existing ingredients item by its ID.

DELETE /ingredients/:id
Delete a specific ingredients item by its ID.
```

# Drinks REST API

```sh
# list of all drinks items
GET /drinks
Retrieve a list of all ingredients items.

POST /drinks
Create a new drinks item.

GET /drinks/:id
Retrieve a specific drinks item by its ID.

PUT /drinks/:id
Update an existing drinks item by its ID.

DELETE /drinks/:id
Delete a specific drinks item by its ID.
```

# Review REST API

```sh
# list of all review items
GET /review
Retrieve a list of all review items.

POST /review
Create a new review item.

GET /review/:id
Retrieve a specific review item by its ID.

PUT /review/:id
Update an existing review item by its ID.

DELETE /review/:id
Delete a specific review item by its ID.
```






# DB Structure

```sh
// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs


Table dish{
    id          bigserial   [primary key]
    createdAt   timestamp(0) 
    updatedAt   timestamp(0) 
    name        text                       
    description text                    
    price       numeric(10, 2)            
}

Table ingredients {
    id         bigserial  [primary key]
    createdAt  timestamp(0) 
    updatedAt  timestamp(0) 
    name       text                       
    quantity   integer                  
    dish_id    bigserial 
}


Table drinks {

    id          bigserial [primary key]
    createdAt   timestamp(0) 
    updatedAt   timestamp(0) 
    name        text                       
    description text                      
    price       numeric(10, 2)            
}

Table review {
    id                        bigserial 
    dish_id                   INTEGER 
    drink_id                  INTEGER 
    rating                    INT 
    comment                   TEXT 
    createdAt                 timestamp(0) 
    updatedAt                 timestamp(0) 
}


Ref: "dish"."id" < "ingredients"."dish_id"
Ref: "dish"."id" < "review"."dish_id"
Ref: "drinks"."id" < "review"."drink_id"
}
```

# Web site where you can see the project

```sh
https://dishes-api.onrender.com/api/v1/dishes
```
