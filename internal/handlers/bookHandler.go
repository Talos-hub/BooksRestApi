package handlers

import "github.com/Talos-hub/BooksRestApi/internal/adapters"

type bookHandler struct {
	logger  adapters.Logger  // needs for write logs
	storage adapters.Storage // is is a storage behavior,
	// it might be Postgresqls, SQLLite and stuff
}
