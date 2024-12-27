package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "errors"
    "strings"
    "github.com/yuvakkrishnan/backend/logger"
)

func DecryptAES(encryptedData []byte, key []byte, log *logger.Logger) ([]byte, error) {
    log.Info("Starting decryption process")

    log.Infof("Received raw key: %s", string(key))
    processedKey := strings.ReplaceAll(string(key), "-", "+")
    processedKey = strings.ReplaceAll(processedKey, "_", "/")

    // Add padding if required
    if len(processedKey)%4 != 0 {
        processedKey = processedKey + strings.Repeat("=", 4-(len(processedKey)%4))
    }

    log.Infof("Processed key: %s", processedKey)
    decoded, err := base64.StdEncoding.DecodeString(processedKey)
    if err != nil {
        log.Errorf("Failed to decode encryption key: %v", err)
        return nil, err
    }

    block, err := aes.NewCipher(decoded)
    if err != nil {
        log.Errorf("Error creating cipher block: %v", err)
        return nil, err
    }

    if len(encryptedData) < aes.BlockSize {
        log.Info("Encrypted data is too short")
        return nil, errors.New("encrypted data is too short")
    }

    iv := encryptedData[:aes.BlockSize]
    encryptedData = encryptedData[aes.BlockSize:]

    if len(encryptedData)%aes.BlockSize != 0 {
        log.Info("Encrypted data is not a multiple of the block size")
        return nil, errors.New("encrypted data is not a multiple of the block size")
    }

    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(encryptedData, encryptedData)

    // Remove padding
    padding := int(encryptedData[len(encryptedData)-1])
    if padding > aes.BlockSize || padding == 0 {
        log.Info("Invalid padding")
        return nil, errors.New("invalid padding")
    }

    for i := 0; i < padding; i++ {
        if encryptedData[len(encryptedData)-1-i] != byte(padding) {
            log.Info("Invalid padding")
            return nil, errors.New("invalid padding")
        }
    }

    decryptedData := encryptedData[:len(encryptedData)-padding]
    log.Info("Decryption process completed successfully")
    return decryptedData, nil
}

