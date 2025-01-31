package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"

	"github.com/yuvakkrishnan/backend/logger"
)

func DecryptAES(encryptedData []byte, key []byte, logger *logger.Logger) ([]byte, error) {
	// Key length validation
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		logger.Error("Invalid key length. Expected 16, 24, or 32 bytes.")
		return nil, errors.New("invalid key length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error("Failed to create cipher block: " + err.Error())
		return nil, err
	}

	if len(encryptedData) < aes.BlockSize {
		logger.Error("Encrypted data is too short")
		return nil, errors.New("encrypted data is too short")
	}

	iv := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptedData, encryptedData)

	// Remove padding
	padding := int(encryptedData[len(encryptedData)-1])
	if padding < 1 || padding > aes.BlockSize {
		logger.Error("Invalid padding: padding value out of range")
		return nil, errors.New("invalid padding or corrupted data")
	}

	for i := 0; i < padding; i++ {
		if encryptedData[len(encryptedData)-1-i] != byte(padding) {
			logger.Error("Invalid padding: padding bytes do not match")
			return nil, errors.New("invalid padding or corrupted data")
		}
	}

	decryptedData := encryptedData[:len(encryptedData)-padding]
	logger.Info("Decryption successful")
	return decryptedData, nil
}
