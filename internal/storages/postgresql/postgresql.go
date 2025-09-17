package postgresql

import (
	"context"
	"fmt"

	"github.com/Talos-hub/BooksRestApi/internal/models"
	"github.com/Talos-hub/BooksRestApi/internal/storages/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool   *pgxpool.Pool
	config *config.DatabaseConfig
}

// NewPosgresStorage create new PostgresStorage that implemented Storage interface
func NewPostgresStorage(config *config.DatabaseConfig) (*PostgresStorage, error) {
	pgxconf, err := pgxpool.ParseConfig(config.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("faild to parse connection string: %w", err)
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
		return nil, fmt.Errorf("faild to create connection pool: %w", err)
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

func (p *PostgresStorage) GetAll() ([]models.Book, error) {
	query := `
	SELECET
		id,
		title,
		author,
		ganre,
		publication_date
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
	for i := 0; rows.Next(); i++ {
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
			return nil, fmt.Errorf("faild to scan books: %w", err)
		}
		books[i] = book
	}

	return books, nil
}

// Helper function to get total book count
func (p *PostgresStorage) getCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM books`

	var count int
	err := p.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("faild to get books count: %w", err)
	}

	return count, nil

}
