package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/yuvakkrishnan/backend/handlers"
    //"github.com/yuvakkrishnan/backend/utils"
    "github.com/yuvakkrishnan/backend/logger"
    
)

func main() {
    // Initialize logger
    logger := logger.NewLogger()

    // Initialize router
    router := mux.NewRouter()

    // Define routes
    router.HandleFunc("/encrypt", handlers.EncryptHandler(logger)).Methods("POST")
    router.HandleFunc("/decrypt", handlers.DecryptHandler(logger)).Methods("POST")
    router.HandleFunc("/upload", handlers.UploadHandler(logger)).Methods("POST")
    router.HandleFunc("/download", handlers.DownloadHandler(logger)).Methods("GET")

    // Start the server
    logger.Info("Starting the server on port 8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        logger.Error("Server failed to start: " + err.Error())
        log.Fatal(err)
    }
}