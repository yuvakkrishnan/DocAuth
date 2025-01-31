package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yuvakkrishnan/backend/handlers"
	"github.com/yuvakkrishnan/backend/logger"
	"github.com/yuvakkrishnan/backend/middleware"
	//"github.com/yuvakkrishnan/backend/utils"
)

func main() {
	// Initialize logger
	logger := logger.NewLogger()

	// Initialize router
	router := mux.NewRouter()

	// Apply CORS middleware
	router.Use(middleware.CORS)

	// Define routes
	router.HandleFunc("/encrypt", handlers.EncryptHandler(logger)).Methods("POST")
	router.HandleFunc("/decrypt", handlers.DecryptHandler(logger)).Methods("POST")
	router.HandleFunc("/upload", handlers.UploadHandler(logger)).Methods("POST")
	router.HandleFunc("/download/{filename}", handlers.DownloadHandler(logger)).Methods("GET")

	logger.Info("Starting the server on port 8080")
	if err := http.ListenAndServe(":9090", router); err != nil {
		logger.Error("Server failed to start: " + err.Error())
		log.Fatal(err)

	}
}
