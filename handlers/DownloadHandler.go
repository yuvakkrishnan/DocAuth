package handlers

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yuvakkrishnan/backend/logger"
	"github.com/yuvakkrishnan/backend/utils"
)

func DownloadHandler(logger *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use mux.Vars to get path parameters
		vars := mux.Vars(r)
		filename := vars["filename"]

		key := r.URL.Query().Get("key")
		if filename == "" || key == "" {
			logger.Error("Missing filename or encryption key")
			http.Error(w, "Missing filename or encryption key", http.StatusBadRequest)
			return
		}

		// Remove all whitespace characters (spaces, tabs, newlines, etc.)
		key = strings.Map(func(r rune) rune {
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
				return -1
			}
			return r
		}, key)

		logger.Info("Sanitized key: " + key)
		logger.Infof("Key length: %d", len(key))

		// Check if the key is a valid base64 string
		if !utils.IsValidBase64(key) {
			logger.Error("Invalid encryption key format")
			http.Error(w, "Invalid encryption key format", http.StatusBadRequest)
			return
		}

		// Decode the key from base64
		keyBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			logger.Error("Failed to decode encryption key: " + err.Error())
			http.Error(w, "Invalid encryption key", http.StatusBadRequest)
			return
		}

		encryptedContent, err := downloadFromS3(filename)
		if err != nil {
			logger.Error("Failed to download file from S3: " + err.Error())
			http.Error(w, "Failed to download file", http.StatusInternalServerError)
			return
		}

		decryptedContent, err := utils.DecryptAES(encryptedContent, keyBytes, logger)
		if err != nil {
			logger.Error("Failed to decrypt content: " + err.Error())
			http.Error(w, "Failed to decrypt content", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", "attachment;filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(decryptedContent)
	}
}
