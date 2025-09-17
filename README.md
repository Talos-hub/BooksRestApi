# Books REST API

A well-structured, production-ready RESTful API for managing a collection of books, built in Go.

## Features

- **CRUD Operations:** Create, read, update, and delete books.
- **Clean Architecture:** Separation into Handlers, Services, and Storage layers.
- **Database Agnostic:** Uses interfaces to allow easy swapping of storage backends (currently PostgreSQL).
- **Proper Error Handling:** Custom error types with consistent JSON responses.
- **Validation:** Robust input validation using reflection (custom implementation).
- **Structured Logging:** Uses `slog` for JSON logging to different files per component.
- **Configuration:** Configurable via environment variables.

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL (with `pgx` driver)
- **Routing:** Standard library `net/http`
- **Logging:** Standard library `log/slog`

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL

### Installation

1.  Clone the repo
    ```bash
    git clone https://github.com/your-username/books-rest-api.git
    cd books-rest-api
    ```

2.  Set up your environment variables (see `Configuration` section below).

3.  Run the application
    ```bash
    go run cmd/main.go
    ```

## API Endpoints

| Method | Endpoint      | Description         |
|--------|---------------|---------------------|
| GET    | `/books`      | Get all books       |
| GET    | `/books/{id}` | Get a book by ID    |
| POST   | `/books`      | Create a new book   |
| PUT    | `/books`      | Update a book       |
| DELETE | `/books/{id}` | Delete a book       |
| GET    | `/health`     | Health check        |

## Configuration

The application is configured using environment variables. See `internal/storages/config/config.go` for all options.

Key variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=bookdb
export PORT=:8080
```
## Project structure
```
├── cmd/
│   └── main.go      = Application entry point
├── internal/
│   ├── abstraction/    = Interfaces (Logger, Storage)
│   ├── apperrors/      = Custom error types
│   ├── handlers/       = HTTP handlers
│   ├── models/         = Data models (Book, etc.)
│   ├── services/       = Business logic layer
│   ├── storages/       = Data persistence layer
│   │   ├── config/     = Database configuration
│   │   └── postgresql/ = PostgreSQL implementation
│   └── validations/    = Input validation logic
└── go.mod
```
