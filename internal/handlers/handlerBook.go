package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Talos-hub/BooksRestApi/internal/abstraction"
	"github.com/Talos-hub/BooksRestApi/internal/apperrors"
	"github.com/Talos-hub/BooksRestApi/internal/models"
	"github.com/Talos-hub/BooksRestApi/internal/services"
)

// HandlerBooks is struct that contains methods
// for handle clients requests
// It implemented ServeHTTP
type HandlerBooks struct {
	Service *services.BookService
	logger  abstraction.Logger
}

// TODO
// func ServeHTTP(w, r)
func (h *HandlerBooks) CreateBook(w http.ResponseWriter, r *http.Request) {
	var createdBook models.CreateBookRequest

	err := json.NewDecoder(r.Body).Decode(&createdBook)
	if err != nil {
		h.sendErrorResponse(w, apperrors.NewAppError(400, "invalid JSON", err))
	}
	t := time.Now()

	createdBook.CreatedAt = t

	apperr := h.Service.CreateBook(createdBook)
	if err != nil {
		h.sendErrorResponse(w, apperr)
	}

	h.sendJsonResponse(w, http.StatusCreated, map[string]string{"message": "create new book seccessfully"})
}

// UpdateBook update a book from a storage by id
func (h *HandlerBooks) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook models.UpdateBookRequest
	err := json.NewDecoder(r.Body).Decode(&updateBook)
	if err != nil {
		h.sendErrorResponse(w, apperrors.NewAppError(400, "invalid JSON", err))
	}
	t := time.Now()

	updateBook.UpdatedAt = t
	appErr := h.Service.UpdateBook(updateBook.Book.ID, updateBook)
	if appErr != nil {
		h.sendErrorResponse(w, appErr)
	}

	h.sendJsonResponse(w, http.StatusOK, map[string]string{"message": "update a book successfully"})

}

// GetAllBooks send all books from a storage to a client
func (h *HandlerBooks) GetAllBooks(w http.ResponseWriter) {
	books, err := h.Service.GetBooks()
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}
	h.sendJsonResponse(w, http.StatusOK, books)
}

// GetById send a book by an ID
func (h *HandlerBooks) GetBookById(w http.ResponseWriter, strID string) {
	if len(strID) == 0 {
		h.sendErrorResponse(w, apperrors.NewAppError(400, "invalid book id", errors.New("id cannot be empty")))
		return
	}
	// parse str to uint64
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		h.sendErrorResponse(w, apperrors.NewAppError(400, "invalid book id", errors.New("id cannot be empty")))
		return
	}

	// get a book
	book, appError := h.Service.GetBook(id)
	if appError != nil {
		h.sendErrorResponse(w, appError)
	}

	h.sendJsonResponse(w, http.StatusOK, book)

}

// func CreateBook()
// func UpdateBook()

// DeleteBook delete a book from a storage
// Where is strId, it's an ID of a book
func (h *HandlerBooks) DeleteBook(w http.ResponseWriter, strID string) {
	if len(strID) == 0 {
		h.sendErrorResponse(w, apperrors.NewAppError(400, "invalid book id", errors.New("id cannot be empty")))
		return
	}
	// parse str to uint64
	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		h.sendErrorResponse(w, apperrors.NewAppError(400, "invalid book id", errors.New("id cannot be empty")))
		return
	}

	// if it has an error, send it
	appErr := h.Service.DeleteBook(id)
	if appErr != nil {
		h.sendErrorResponse(w, appErr)
	}

	h.sendJsonResponse(w, http.StatusOK, map[string]string{"message": "Book deleted successfully"})

}

// SendJsonResponse send to client a json response.
// If data is nil it send bad status code
func (h *HandlerBooks) sendJsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-type", "application/json")

	// if data is nil then write a bad status code
	if data == nil {
		h.logger.Error("error send json response data is nil", "data", data)
		w.WriteHeader(http.StatusInternalServerError)

		errorResponse := map[string]string{
			"error": "internal server error",
		}
		err := json.NewEncoder(w).Encode(errorResponse)
		if err != nil {
			h.logger.Error("error send error response", "error", err, "errorResponse", errorResponse["error"])
		}
		return
	}
	// set a header and write status code

	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)

	err := encoder.Encode(data)
	if err != nil {
		h.logger.Error("error encode response", "error", err, "data", data)
	}

}

// sendErrorResponse send to cliend an error, if error is nil,
// it write log and returna
func (h *HandlerBooks) sendErrorResponse(w http.ResponseWriter, appErr *apperrors.AppError) {
	if appErr == nil {
		h.logger.Warn("Error send a app error, error is nil", "appErr", appErr)
		return
	}
	h.logger.Info("API error", "code", appErr.Code, "message", appErr.Message, "error", appErr.Err)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(appErr.Code)

	err := json.NewEncoder(w).Encode(map[string]any{
		"code":    appErr.Code,
		"message": appErr.Message,
	})

	if err != nil {
		h.logger.Error("Error send an erro response", "error", err)
	}

}
