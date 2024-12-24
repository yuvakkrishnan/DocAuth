package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
)

func GenerateAESKey() (string, error) {
    key := make([]byte, 32) // AES-256 key size
    _, err := rand.Read(key)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(key), nil
}

func EncryptAES(content []byte, key string) ([]byte, error) {
    decodedKey, err := base64.StdEncoding.DecodeString(key)
    if err != nil {
        return nil, err
    }
    block, err := aes.NewCipher(decodedKey)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, 12) // GCM standard nonce size
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    ciphertext := aesGCM.Seal(nonce, nonce, content, nil)
    return ciphertext, nil
}
func DecryptAES(ciphertext []byte, key string) ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

