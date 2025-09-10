package models

import "time"

// Book is model that implemented behavior a real book.
type Book struct {
	ID              uint64    `json:"id"`              // unique Id
	Title           string    `json:"title"`           // Name of a Book
	Ganre           string    `json:"ganre"`           // for example Adnveture, Roman
	PublicationData time.Time `json:"publicationData"` // for instance 1970
	Author          string    `json:"author"`          // full name of Author
	CreatedAt       time.Time `json:"createdAt"`       // time when is was created
	UpdatedAt       time.Time `json:"updateAt"`        // time when is was updated
}
