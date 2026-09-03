package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func getEncryptionKey() ([]byte, error) {
	keyString := os.Getenv("ENCRYPTION_KEY")

	if keyString == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY is not set")
	}

	key, err := base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		return nil, fmt.Errorf("invalid ENCRYPTION_KEY: %w", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("ENCRYPTION_KEY must be exactly 32 bytes")
	}

	return key, nil
}

func Encrypt(plaintext string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	encrypted := gcm.Seal(
		nonce,
		nonce,
		[]byte(plaintext),
		nil,
	)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func Decrypt(ciphertext string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return "", fmt.Errorf("invalid encrypted data")
	}

	nonce := data[:nonceSize]
	encryptedData := data[nonceSize:]

	plaintext, err := gcm.Open(
		nil,
		nonce,
		encryptedData,
		nil,
	)

	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %w", err)
	}

	return string(plaintext), nil
}
