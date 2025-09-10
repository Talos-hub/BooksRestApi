package models

type Book struct {
	ID              uint64 // unique Id
	Book            string // Name of a Book
	Ganre           string // for example Adnveture, Roman
	PublicationData uint64 // for instance 1970
	Author          string // full name of Author
}
