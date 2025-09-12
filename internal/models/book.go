package models

import (
	"time"
)

type GeneralBook struct {
	ID              uint64    `json:"id"`              // unique Id
	Title           string    `json:"title"`           // Name of a Book
	Genre           string    `json:"ganre"`           // for example Adnveture, Roman
	PublicationDate time.Time `json:"publicationData"` // for instance 1970
	Author          string    `json:"author"`
}

// Book is model that implemented behavior a real book.
type Book struct {
	General   GeneralBook `json:"general"`   // it is simple book
	CreatedAt time.Time   `json:"createdAt"` // time when is was created
	UpdatedAt time.Time   `json:"updateAt"`  // time when is was updated
}

type UpdateBookRequest struct {
	Book GeneralBook
}

type CreateBookRequest struct {
	Book GeneralBook
}
