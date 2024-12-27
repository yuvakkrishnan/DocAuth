package handlers

import (
        "encoding/base64"
        "github.com/yuvakkrishnan/backend/utils"
        "github.com/yuvakkrishnan/backend/logger"
        "io"
        "net/http"
        "strings"
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

                // Trim spaces and ensure consistent key format
                key = strings.TrimSpace(key)

                // Decode the key (assuming it's already base64 encoded)
                keyBytes, err := base64.StdEncoding.DecodeString(key) 
                if err != nil {
                        logger.Error("Invalid base64 encoding for key: " + err.Error())
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