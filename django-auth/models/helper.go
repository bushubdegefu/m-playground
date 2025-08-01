package models

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/bushubdegefu/m-playground/configs"
)

// Helper for pagination
type Pagination struct {
	Page   int
	Size   int
	Offset int
}

// Combine password and salt then hash them using the SHA-512
func HashFunc(password string) string {

	// var salt []byte
	// get salt from env variable
	salt := []byte(configs.AppConfig.Get("SECRETE_SALT"))

	// Convert password string to byte slice
	var pwdByte = []byte(password)

	// Create sha-512 hasher
	var sha512 = sha512.New()

	pwdByte = append(pwdByte, salt...)

	sha512.Write(pwdByte)

	// Get the SHA-512 hashed password
	var hashedPassword = sha512.Sum(nil)

	// Convert the hashed to hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPassword)
	return hashedPasswordHex
}
