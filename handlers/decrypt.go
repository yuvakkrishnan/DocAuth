package handlers

import (
	"github.com/yuvakkrishnan/backend/utils"
	"io"
	"net/http"
)

func DecryptHandler(logger *utils.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logger.Info("Handling file decryption")

        file, _, err := r.FormFile("file")
        if err != nil {
            logger.Error("Failed to retrieve file from request: ")
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

        encryptedContent, err := io.ReadAll(file)
        if err != nil {
            logger.Error("Failed to read encrypted file content: ")
            http.Error(w, "Failed to read encrypted file content", http.StatusInternalServerError)
            return
        }

        decryptedContent, err := utils.DecryptAES(encryptedContent, key)
        if err != nil {
            logger.Error("Failed to decrypt content: ")
            http.Error(w, "Failed to decrypt content", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/pdf")
        w.Write(decryptedContent)
        logger.Info("File decrypted successfully")
    }
}

