package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dollarkillerx/unregistry/internal/handler"
	"github.com/dollarkillerx/unregistry/internal/storage"
	"github.com/dollarkillerx/unregistry/pkg/auth"
	"github.com/gorilla/mux"
)

func main() {
	// Get configuration from environment variables
	token := os.Getenv("TOKEN")
	if token == "" {
		token = "default-token"
		log.Println("Warning: Using default token. Set TOKEN environment variable for production.")
	}

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "/data"
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = "0.0.0.0:8080"
	}

	// Initialize storage
	storage := storage.New(dataPath)

	// Initialize handlers
	fileHandler := handler.NewFileHandler(storage)
	imageHandler := handler.NewImageHandler(storage)

	// Setup router
	r := mux.NewRouter()

	// API routes with authentication
	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.AuthMiddleware(token))

	// File routes
	api.HandleFunc("/file/upload", fileHandler.Upload).Methods("POST")
	api.HandleFunc("/file/download/{filename}", fileHandler.Download).Methods("GET")
	api.HandleFunc("/file/list", fileHandler.List).Methods("GET")
	api.HandleFunc("/file/{filename}", fileHandler.Delete).Methods("DELETE")

	// Image routes
	api.HandleFunc("/img/upload", imageHandler.Upload).Methods("POST")
	api.HandleFunc("/img/download/{name}", imageHandler.Download).Methods("GET")
	api.HandleFunc("/img/list", imageHandler.List).Methods("GET")
	api.HandleFunc("/img/{name}", imageHandler.Delete).Methods("DELETE")

	// Health check endpoint (no auth required)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	log.Printf("Server starting on %s", listenAddr)
	log.Printf("Data path: %s", dataPath)
	
	if err := http.ListenAndServe(listenAddr, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}