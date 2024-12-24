package tests

import (
	"backend/utils"
	"testing"
)

func TestGenerateAESKey(t *testing.T) {
	key, err := utils.GenerateAESKey()
	if err != nil {
		t.Fatalf("Failed to generate AES key: %v", err)
	}
	if len(key) == 0 {
		t.Fatalf("Generated key is empty")
	}
}

func TestEncryptAndDecryptAES(t *testing.T) {
	key, _ := utils.GenerateAESKey()
	content := []byte("test content")

	encrypted, err := utils.EncryptAES(content, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := utils.DecryptAES(encrypted, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if string(decrypted) != string(content) {
		t.Fatalf("Decrypted content does not match original")
	}
}
