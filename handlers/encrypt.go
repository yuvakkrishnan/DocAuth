package handlers

import (
    "encoding/base64"
    "net/http"
    "io"
    "strings"
	"encoding/json"
    "github.com/yuvakkrishnan/backend/utils"
	"github.com/yuvakkrishnan/backend/logger"
)

func EncryptHandler(logger *logger.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        file, header, err := r.FormFile("file")
        if err != nil {
            logger.Error("Failed to retrieve file from request")
            http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
            return
        }
        defer file.Close()

        fileContent, err := io.ReadAll(file)
        if err != nil {
            logger.Error("Failed to read file content")
            http.Error(w, "Failed to read file content", http.StatusInternalServerError)
            return
        }

        key, err := utils.GenerateAESKey()
        if err != nil {
            logger.Error("Failed to generate encryption key")
            http.Error(w, "Failed to generate encryption key", http.StatusInternalServerError)
            return
        }

        // Ensure the key is clean (valid Base64)
        key = strings.TrimSpace(key)
        key = base64.StdEncoding.EncodeToString([]byte(key))

        keyBytes, err := base64.StdEncoding.DecodeString(key)
        if err != nil {
            logger.Error("Failed to decode encryption key: " + err.Error())
            http.Error(w, "Invalid encryption key", http.StatusBadRequest)
            return
        }

        encryptedContent, err := utils.EncryptAES(fileContent, keyBytes)
        if err != nil {
            logger.Error("Failed to encrypt file content")
            http.Error(w, "Failed to encrypt file", http.StatusInternalServerError)
            return
        }

        filename := header.Filename
        if err := uploadToS3(filename, encryptedContent); err != nil {
            logger.Error("Failed to upload encrypted file to S3")
            http.Error(w, "Failed to upload file to storage", http.StatusInternalServerError)
            return
        }

        response := map[string]string{
            "message":      "File uploaded and encrypted successfully",
            "filename":     filename,
            "encryptionKey": key,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        logger.Info("File uploaded and encrypted successfully")
    }
}
