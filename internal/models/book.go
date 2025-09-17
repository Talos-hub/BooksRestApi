package models

import (
	"time"
)

type GeneralBook struct {
	ID              uint64    `json:"id" db:"id"`                            // unique Id
	Title           string    `json:"title" db:"title"`                      // Name of a Book
	Genre           string    `json:"genre" db:"genre"`                      // for example Adnveture, Roman
	PublicationDate time.Time `json:"publicationDate" db:"publication_date"` // for instance 1970
	Author          string    `json:"author" db:"author"`
}

// Book is model that implemented behavior a real book.
type Book struct {
	General   GeneralBook `json:"general"`                   // it is simple book
	CreatedAt time.Time   `json:"createdAt" db:"created_at"` // time when is was created
	UpdatedAt time.Time   `json:"updateAt" db:"updated_at"`  // time when is was updated
}

type UpdateBookRequest struct {
	Book      GeneralBook `jsons:"book"`
	UpdatedAt time.Time   `json:"-"` // time when is was updated
}

type CreateBookRequest struct {
	Book      GeneralBook `json:"book"`
	CreatedAt time.Time   `json:"-"` // time when is was created
}
