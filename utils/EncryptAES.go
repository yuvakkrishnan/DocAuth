package utils

import (
    "crypto/aes"
    "crypto/rand"
	"bytes"
	"crypto/cipher"

)

func EncryptAES(content []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    // Add padding
    padding := aes.BlockSize - len(content)%aes.BlockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    content = append(content, padtext...)

    ciphertext := make([]byte, aes.BlockSize+len(content))
    iv := ciphertext[:aes.BlockSize]
    if _, err := rand.Read(iv); err != nil {
        return nil, err
    }

    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(ciphertext[aes.BlockSize:], content)

    return ciphertext, nil
}
