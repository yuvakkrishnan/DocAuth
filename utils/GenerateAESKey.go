package utils

import (
    "crypto/rand"
    "encoding/base64"
)

func GenerateAESKey() (string, error) {
    key := make([]byte, 32) // AES-256 key size (32 bytes)
    if _, err := rand.Read(key); err != nil {
        return "", err
    }

    // Encode the key to base64
    encodedKey := base64.StdEncoding.EncodeToString(key)
    return encodedKey, nil
}
