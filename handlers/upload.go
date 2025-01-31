package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/yuvakkrishnan/backend/logger"
	"github.com/yuvakkrishnan/backend/utils"
)

func UploadHandler(logger *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received file upload request")

		// Allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		logger.Info("Handling file upload")

		// Parse file from the request
		file, header, err := r.FormFile("file")
		if err != nil {
			logger.Error("Failed to retrieve file from request: " + err.Error())
			http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			logger.Error("Failed to read file content: " + err.Error())
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}

		// Generate AES encryption key
		key, err := utils.GenerateAESKey()
		if err != nil {
			logger.Error("Failed to generate encryption key: " + err.Error())
			http.Error(w, "Failed to generate encryption key", http.StatusInternalServerError)
			return
		}

		// Convert key from base64 string to byte slice
		keyBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			logger.Error("Failed to decode encryption key: " + err.Error())
			http.Error(w, "Invalid encryption key", http.StatusBadRequest)
			return
		}

		// Encrypt the file content
		encryptedContent, err := utils.EncryptAES(fileContent, keyBytes)
		if err != nil {
			logger.Error("Failed to encrypt file content: " + err.Error())
			http.Error(w, "Failed to encrypt file content", http.StatusInternalServerError)
			return
		}

		// Store the encrypted file to S3
		filename := header.Filename
		if err := uploadToS3(filename, encryptedContent); err != nil {
			logger.Error("Failed to upload encrypted file to S3: " + err.Error())
			http.Error(w, "Failed to upload file to storage", http.StatusInternalServerError)
			return
		}

		// Prepare response
		response := map[string]string{
			"message":       "File uploaded and encrypted successfully",
			"filename":      filename,
			"encryptionKey": key, // return the key as a string
		}

		// Respond with JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		logger.Info("File uploaded and encrypted successfully")
	}
}
