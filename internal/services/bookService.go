package services

import (
	"errors"

	"github.com/Talos-hub/BooksRestApi/internal/abstraction"
	"github.com/Talos-hub/BooksRestApi/internal/apperrors"
	"github.com/Talos-hub/BooksRestApi/internal/models"
	"github.com/Talos-hub/BooksRestApi/internal/validations"
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
		s.logger.Info("Failed to get book by ID", "id", id, "error", err)
		return models.Book{}, apperrors.NewAppError(404, "book not found", err)
	}

	return book, nil
}

// Created created new book and save it to storage
func (s *bookService) CreateBook(book models.CreateBookRequest) *apperrors.AppError {
	// validation
	err := validations.Validate(book)
	if err != nil {
		// if someone use it worng it returns ValidationReflectErr
		// For instance: if parameter is func it returns the error
		if errors.Is(err, &apperrors.ValidationReflectErr{}) {
			s.logger.Error("Error validation", "error", err)
			return apperrors.NewAppError(500, "error creating a book", err)
		}
		return apperrors.NewAppError(400, "invalid book data", err)
	}

	// created new book
	newBook := models.Book{
		General:   book.Book,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.CreatedAt,
	}
	// save a book
	err = s.storage.Save(newBook)
	if err != nil {
		s.logger.Error("Error save a book", "error", err)
		return apperrors.NewAppError(500, "faild to create a book", err)
	}

	return nil
}

func (s *bookService) UpdateBook(id uint64, update models.UpdateBookRequest) *apperrors.AppError {
	book, err := s.storage.GetById(id)
	if err != nil {
		s.logger.Info("faild to update a book")
		return apperrors.NewAppError(404, "a book not found", err)
	}

	// validation
	err = validations.Validate(update)
	if err != nil {
		// if someone use it worng it returns ValidationReflectErr
		// For instance: if parameter is func it returns the error
		if errors.Is(err, &apperrors.ValidationReflectErr{}) {
			s.logger.Error("Error validation", "error", err)
			return apperrors.NewAppError(500, "error update a book", err)
		}
		return apperrors.NewAppError(400, "invalid book data", err)
	}

}
