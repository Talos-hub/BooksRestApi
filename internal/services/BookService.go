package services

import (
	"github.com/Talos-hub/BooksRestApi/internal/abstraction"
)

// BookService impemented handlers for hadle book.
// It contains logger and storage interface
type bookService struct {
	logger  abstraction.Logger  // Logger that needed for write logs
	storage abstraction.Storage // it might be any strage for instance Postgresqls, Sqlite, json
}

func NewBookService(logger abstraction.Logger, storage abstraction.Storage) *bookService {
	return &bookService{
		logger:  logger,
		storage: storage,
	}
}
