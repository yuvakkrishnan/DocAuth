package handlers

import (
	"github.com/yuvakkrishnan/backend/utils"
	"encoding/json"
	"io"
	"net/http"
)

func EncryptHandler(logger *utils.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse file from the request
		file, header, err := r.FormFile("file")
		if err != nil {
			logger.Error("Failed to retrieve file from request")
			http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			logger.Error("Failed to read file content")
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}

		// Generate AES encryption key
		key, err := utils.GenerateAESKey()
		if err != nil {
			logger.Error("Failed to generate encryption key")
			http.Error(w, "Failed to generate encryption key", http.StatusInternalServerError)
			return
		}

		// Encrypt the file content
		encryptedContent, err := utils.EncryptAES(fileContent, key)
		if err != nil {
			logger.Error("Failed to encrypt file content")
			http.Error(w, "Failed to encrypt file", http.StatusInternalServerError)
			return
		}

		// Store the encrypted file to S3
		filename := header.Filename
		if err := uploadToS3(filename, encryptedContent); err != nil {
			logger.Error("Failed to upload encrypted file to S3")
			http.Error(w, "Failed to upload file to storage", http.StatusInternalServerError)
			return
		}

		// Prepare response
		response := map[string]string{
			"message":      "File uploaded and encrypted successfully",
			"filename":     filename,
			"encryptionKey": key,
		}

		// Respond with JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		logger.Info("File uploaded and encrypted successfully")
	}
}
