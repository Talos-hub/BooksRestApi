package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Talos-hub/BooksRestApi/internal/abstraction"
	"github.com/Talos-hub/BooksRestApi/internal/apperrors"
	"github.com/Talos-hub/BooksRestApi/internal/services"
)

// HandlerBooks is struct that contains methods
// for handle clients requests
// It implemented ServeHTTP
type HandlerBooks struct {
	Service *services.BookService
	logger  abstraction.Logger
}

//TODO
// func ServeHTTP(w, r)

//TODO
// func GetAll()
// func GetByID()
// func CreateBook()
// func UpdateBook()
// func DeleteBook()

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
// it write log and return
func (h *HandlerBooks) sendErrorResponse(w http.ResponseWriter, appErr *apperrors.AppError) {
	if appErr == nil {
		h.logger.Warn("Error send a app error, error is nil", "appErr", appErr)
		return
	}
	h.logger.Error("API error", "code", appErr.Code, "message", appErr.Message, "error", appErr.Err)

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
