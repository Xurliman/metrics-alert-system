package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"sync"
)

var (
	rsaPublicKey *rsa.PublicKey // Store the public key
	once         sync.Once
)

// LoadPublicKey reads and parses the RSA public key once during startup
func LoadPublicKey(cryptoKey string) error {
	pubKeyData, err := os.ReadFile(cryptoKey)
	if err != nil {
		return fmt.Errorf("failed to read public key: %v", err)
	}

	block, _ := pem.Decode(pubKeyData)
	if block == nil {
		return fmt.Errorf("invalid PEM format")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}

	key, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("invalid public key")
	}

	rsaPublicKey = key
	return nil
}

func Encrypt(plaintext []byte) ([]byte, error) {
	if rsaPublicKey == nil {
		return plaintext, nil
	}

	// Encrypt using RSA-OAEP
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKey, plaintext, nil)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %v", err)
	}

	return encryptedData, nil
}
