package middlewares

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/server/utils"
	"github.com/Xurliman/metrics-alert-system/internal/log"
	"go.uber.org/zap"
	"io"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type DecryptingMiddleware struct{}

var (
	rsaPrivateKey *rsa.PrivateKey // Store the private key
	once          sync.Once
)

func NewDecryptingMiddleware(cryptoKey string) Middleware {
	once.Do(func() {
		if err := LoadPrivateKey(cryptoKey); err != nil {
			log.Fatal("failed to load public key", zap.Error(err))
		}
	})
	return &DecryptingMiddleware{}
}

// Handle decrypts incoming encrypted request data
func (h DecryptingMiddleware) Handle(ctx *gin.Context) {
	if rsaPrivateKey == nil {
		ctx.Next()
		return
	}

	// Read the encrypted data from the request body
	encryptedData, err := ctx.GetRawData()
	if err != nil {
		utils.JSONError(ctx, fmt.Errorf("failed to read request body"))
		return
	}

	// Decrypt using RSA-OAEP (Optimal Asymmetric Encryption Padding)
	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, encryptedData, nil)
	if err != nil {
		utils.JSONError(ctx, fmt.Errorf("decryption failed: %v", err))
		return
	}

	// Replace the body with the decrypted data
	ctx.Request.Body = io.NopCloser(bytes.NewReader(decryptedData))

	// Allow request to proceed
	ctx.Next()
}

// LoadPrivateKey reads and parses the RSA private key once during startup
func LoadPrivateKey(cryptoKey string) error {
	privKeyData, err := os.ReadFile(cryptoKey)
	if err != nil {
		return fmt.Errorf("failed to read private key: %v", err)
	}

	block, _ := pem.Decode(privKeyData)
	if block == nil {
		return fmt.Errorf("invalid PEM format")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	rsaPrivateKey = privKey
	return nil
}
