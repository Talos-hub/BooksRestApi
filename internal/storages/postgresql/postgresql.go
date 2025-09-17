package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Talos-hub/BooksRestApi/internal/abstraction"
	"github.com/Talos-hub/BooksRestApi/internal/models"
	"github.com/Talos-hub/BooksRestApi/internal/storages/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool   *pgxpool.Pool
	config *config.DatabaseConfig
	logger abstraction.Logger
}

// NewPosgresStorage create new PostgresStorage that implemented Storage interface
func NewPostgresStorage(config *config.DatabaseConfig, logger abstraction.Logger) (*PostgresStorage, error) {
	pgxconf, err := pgxpool.ParseConfig(config.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	//set connection pool settings
	pgxconf.MaxConns = int32(config.MaxConns)
	pgxconf.MinConns = int32(config.MinConns)
	pgxconf.MaxConnLifetime = config.ConnMaxLifeTime
	pgxconf.MaxConnIdleTime = config.ConnMaxIdleTime
	pgxconf.HealthCheckPeriod = config.HealthCheckPeriod

	// create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxconf)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	// ping to database
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("faild to ping database: %w", err)
	}

	// create a table book if it not exist
	err = initTable(config, pool)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		pool:   pool,
		config: config,
		logger: logger,
	}, nil
}

// initTable create a table if it not exist
func initTable(config *config.DatabaseConfig, pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		author VARCHAR(100) NOT NULL,
		genre VARCHAR(100) NOT NULL,
		publication_date TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	`
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("faild to init database table: %w", err)
	}

	return nil

}

// GetAll return all books from storage
func (p *PostgresStorage) GetAll() ([]models.Book, error) {
	query := `
	SELECT
		id,
		title,
		author,
		genre,
		publication_date,
		created_at,
		updated_at
	FROM books
	ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()

	//get amount of books
	count, err := p.getCount(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		p.logger.Error("Faild to query books", "error", err)
		return nil, fmt.Errorf("faild to query books: %w", err)
	}

	defer rows.Close()

	// created books avoid allocations
	var books []models.Book = make([]models.Book, count)
	var book models.Book

	// it uses books[i] = book insted append because it more performance
	// Direct assignment - very fast :
	// No memory allocations during the loop,
	// Direct memory access - O(1) time complexity
	i := 0
	for rows.Next() {
		err := rows.Scan(
			&book.General.ID,
			&book.General.Title,
			&book.General.Author,
			&book.General.Genre,
			&book.General.PublicationDate,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			p.logger.Error("Faild to scan books", "error", err)
			return nil, fmt.Errorf("faild to scan books: %w", err)
		}
		if i < count {
			books[i] = book
			i++
		}
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		p.logger.Error("Error iterating rows", "error", err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return books, nil
}

// GetById return a book by id
func (p *PostgresStorage) GetById(id uint64) (models.Book, error) {
	query := `
	SELECT
		id,
		title,
		author,
		genre,
		publication_date,
		created_at,
		updated_at
	FROM books
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()

	var book models.Book
	err := p.pool.QueryRow(ctx, query, id).Scan(
		&book.General.ID,
		&book.General.Title,
		&book.General.Author,
		&book.General.Genre,
		&book.General.PublicationDate,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Book{}, fmt.Errorf("book with id %d not found", id)
		}
		p.logger.Error("Faild to get book", "error", err)
		return models.Book{}, fmt.Errorf("failed to get book: %w", err)
	}
	return book, nil
}

// Save add a book to database
func (p *PostgresStorage) Save(book models.Book) error {
	query := `
	INSERT INTO books (title, author, genre, publication_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()

	err := p.pool.QueryRow(ctx, query,
		book.General.Title,
		book.General.Author,
		book.General.Genre,
		book.General.PublicationDate,
		book.CreatedAt,
		book.UpdatedAt,
	).Scan(&book.General.ID)

	if err != nil {
		p.logger.Error("Failed to save book", "error", err)
		return fmt.Errorf("failed to save book: %w", err)
	}

	return nil

}

// Update update a book into database
func (p *PostgresStorage) Update(book models.Book) error {
	query := `
	UPDATE books 
	SET 
		title = $1, 
		author = $2, 
		genre = $3, 
		publication_date = $4, 
		updated_at = $5
	WHERE id = $6
	`

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()

	result, err := p.pool.Exec(ctx, query, book.General.Title,
		book.General.Author,
		book.General.Genre,
		book.General.PublicationDate,
		book.UpdatedAt,
		book.General.ID,
	)

	if err != nil {
		p.logger.Error("Failed to update book", "error", err)
		return fmt.Errorf("failed to update book: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("book with id: %d not found", book.General.ID)
	}

	return nil

}

// Delete delete a book by id
func (p *PostgresStorage) Delete(id uint64) error {
	query := `
	DELETE FROM books
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()

	result, err := p.pool.Exec(ctx, query, id)
	if err != nil {
		p.logger.Error("Failed to delete a book", "error", err)
		return fmt.Errorf("failed to delete a book: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("book with id: %d not found", id)
	}
	return nil
}

// Close close a storage
func (p *PostgresStorage) Close() error {
	p.pool.Close()
	return nil
}

// Helper function to get total book count
func (p *PostgresStorage) getCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM books`

	var count int
	err := p.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		p.logger.Error("Failed to get book id", "error", err)
		return 0, fmt.Errorf("failed to get books count: %w", err)
	}

	return count, nil

}
