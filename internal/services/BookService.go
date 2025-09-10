package services

import (
	"github.com/Talos-hub/BooksRestApi/internal/abstraction"
)

type BookService struct {
	logger  abstraction.Logger
	storage abstraction.Storage
}
