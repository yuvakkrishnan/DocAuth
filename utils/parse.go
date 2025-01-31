package utils

import (
    "encoding/base64"
    "strings"
)

// IsValidBase64 checks if a string is valid Base64
func IsValidBase64(s string) bool {
    s = strings.TrimSpace(s)

    // Check if Base64 string length is a multiple of 4
    if len(s)%4 != 0 {
        return false
    }

    // Decode to validate
    _, err := base64.StdEncoding.DecodeString(s)
    return err == nil
}
