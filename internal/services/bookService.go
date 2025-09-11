package services

import (
	"github.com/Talos-hub/BooksRestApi/internal/abstraction"
	"github.com/Talos-hub/BooksRestApi/internal/apperrors"
	"github.com/Talos-hub/BooksRestApi/internal/models"
)

// BookService impemented handlers for hadle book.
// It contains logger and storage interface
type bookService struct {
	logger  abstraction.Logger  // Logger that needed for write logs
	storage abstraction.Storage // it might be any strage for instance Postgresqls, Sqlite, json
}

// Construction that set a logger and a storage and returns pointer to bookService
func NewBookService(logger abstraction.Logger, storage abstraction.Storage) *bookService {
	return &bookService{
		logger:  logger,
		storage: storage,
	}
}

// GetBooks returns all books from storage
func (s *bookService) GetBooks() ([]models.Book, *apperrors.AppError) {
	books, err := s.storage.GetAll()
	if err != nil {
		s.logger.Error("Error getting all books", "error", err)
		return nil, apperrors.NewAppError(500, "error getting all books", err)
	}
	return books, nil
}

// GetBook return a book by id
func (s *bookService) GetBook(id uint64) (models.Book, *apperrors.AppError) {
	book, err := s.storage.GetById(id)
	if err != nil {
		s.logger.Error("Failed to get book by ID", "id", id, "error", err)
		return models.Book{}, apperrors.NewAppError(404, "book not found", err)
	}

	return book, nil
}
