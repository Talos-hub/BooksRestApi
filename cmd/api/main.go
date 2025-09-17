package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Talos-hub/BooksRestApi/internal/handlers"
	"github.com/Talos-hub/BooksRestApi/internal/services"
	"github.com/Talos-hub/BooksRestApi/internal/storages/config"
	"github.com/Talos-hub/BooksRestApi/internal/storages/postgresql"
)

func main() {
	// init loggers
	storagelogger := SetLogger("storage_log/storage.log")
	servicelogger := SetLogger("service_log/service.log")
	hanlderslogger := SetLogger("handlers_log/handler.log")

	// database configuration
	conf := config.LoadConfig()
	port := os.Getenv("PORT")

	//database
	storage, err := postgresql.NewPostgresStorage(conf, storagelogger)
	if err != nil {
		log.Printf("Cannot create database connect: %v\n", err)
		log.Fatal(err)
	}

	//Book Service
	bookservice := services.NewBookService(servicelogger, storage)

	//Handler
	handler := handlers.NewHandlerBooks(bookservice, hanlderslogger)

	//new router
	mux := http.NewServeMux()
	// set up routes
	mux.Handle("/books", handler)
	mux.Handle("/books/", handler)
	mux.HandleFunc("/health", healthCheck)

	// create server
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// there are helpers

func SetLogger(path string) *slog.Logger {
	// Extract directory from file path
	dir := filepath.Dir(path)

	// Ensure the directory exists first
	if err := EnsureDirectory(dir); err != nil {
		// Fallback to stdout if directory creation fails
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	// Now create/open the log file
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		// Fallback to stdout if file creation fails
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	// Log to file
	return slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func EnsureDirectory(path string) error {
	cleanPath := filepath.Clean(path)

	// Check if the path exists and is a directory
	info, err := os.Stat(cleanPath)
	if err == nil {
		if info.IsDir() {
			return nil // Directory exists and is valid
		}
		return &os.PathError{Op: "mkdir", Path: cleanPath, Err: os.ErrExist}
	}

	// If path doesn't exist, create it with proper permissions
	if os.IsNotExist(err) {
		if err := os.MkdirAll(cleanPath, 0755); err != nil {
			return err
		}
		return nil
	}

	return err
}
