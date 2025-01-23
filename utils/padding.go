package utils

import (
	"errors"
)

func validatePadding(data []byte, blockSize int) error {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return errors.New("invalid data length for padding")
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return errors.New("padding value out of range")
	}

	for i := 0; i < padding; i++ {
		if data[len(data)-1-i] != byte(padding) {
			return errors.New("padding bytes do not match expected value")
		}
	}
	return nil
}
