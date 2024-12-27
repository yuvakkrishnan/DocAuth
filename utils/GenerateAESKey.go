package utils

import (
    "crypto/rand"
    "encoding/base64"
    "strings"
)

func GenerateAESKey() (string, error) {
    key := make([]byte, 32) // AES-256 key size
    if _, err := rand.Read(key); err != nil {
        return "", err
    }

    // Trim and encode to ensure valid Base64
    keyStr := base64.StdEncoding.EncodeToString(key)
    keyStr = strings.TrimSpace(keyStr)
    return keyStr, nil
}
