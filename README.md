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

# DB Structure

```sh
// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs


Table dish{
    id          bigserial PRIMARY KEY,
    createdAt   timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt   timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name        text                        NOT NULL,
    description text                        NOT NULL,
    price       numeric(10, 2)              NOT NULL
}

Table ingredients {
    id         bigserial PRIMARY KEY,
    createdAt  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name       text                        NOT NULL,
    quantity   integer                     NOT NULL,
    dish_id    bigserial REFERENCES dishes (id)
}
```

# Web site where you can see the project

```sh
https://dishes-api.onrender.com/api/v1/dishes
```
