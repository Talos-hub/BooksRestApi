package adapters

import "github.com/Talos-hub/BooksRestApi/internal/models"

// It  interface that provides Storage.
// It has all methods for work with any storage:
// SqlLite, Postgresql, json file, etc.
type Storage interface {
	GetAll() ([]models.Book, error)         // returns all elements from a storage
	GetById(id uint64) (models.Book, error) // returns one item from a storage by id
	Save(book models.Book) error            // add a book to storage
	Delete(id uint64) error                 // delete a item from storage
	Update(book models.Book) error          // update a item in storage
}
