package handlers

import (
    "encoding/base64"
    "net/http"
    "net/url"
    "strings"
    "github.com/yuvakkrishnan/backend/utils"
    "github.com/yuvakkrishnan/backend/logger"
)

func DownloadHandler(logger *logger.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        filename := r.URL.Query().Get("filename")
        key := r.URL.Query().Get("key")

        if filename == "" || key == "" {
            logger.Error("Missing filename or encryption key")
            http.Error(w, "Missing filename or encryption key", http.StatusBadRequest)
            return
        }

        // Trim spaces and unescape URL-encoded characters
        key = strings.TrimSpace(key)
        key, err := url.QueryUnescape(key)
        if err != nil {
            logger.Error("Failed to unescape encryption key: " + err.Error())
            http.Error(w, "Failed to unescape encryption key", http.StatusBadRequest)
            return
        }

        logger.Info("Received raw key: " + key)

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

        // Write the decrypted content to the response
        w.Header().Set("Content-Disposition", "attachment; filename="+filename)
        w.Header().Set("Content-Type", "application/octet-stream")
        w.Write(decryptedContent)
    }
}