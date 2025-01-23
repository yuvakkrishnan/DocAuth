package handlers

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/yuvakkrishnan/backend/logger"
	"github.com/yuvakkrishnan/backend/utils"
)

func DecryptHandler(logger *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Handling file decryption")

		file, _, err := r.FormFile("file")
		if err != nil {
			logger.Error("Failed to retrieve file from request: " + err.Error())
			http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		key := r.FormValue("key")
		if key == "" {
			logger.Error("Missing encryption key")
			http.Error(w, "Missing encryption key", http.StatusBadRequest)
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

		// Decode the key from base64
		keyBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			logger.Error("Failed to decode encryption key: " + err.Error())
			http.Error(w, "Invalid encryption key", http.StatusBadRequest)
			return
		}

		encryptedContent, err := io.ReadAll(file)
		if err != nil {
			logger.Error("Failed to read encrypted file content: " + err.Error())
			http.Error(w, "Failed to read encrypted file content", http.StatusInternalServerError)
			return
		}

		decryptedContent, err := utils.DecryptAES(encryptedContent, keyBytes, logger)
		if err != nil {
			logger.Error("Failed to decrypt content: " + err.Error())
			http.Error(w, "Failed to decrypt content", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=decrypted_file")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(decryptedContent)
	}
}
