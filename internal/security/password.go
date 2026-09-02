package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
	saltLength   = 16
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)

	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory,
		argonTime,
		argonThreads,
		saltEncoded,
		hashEncoded,
	), nil
}

func VerifyPassword(password string, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return false, fmt.Errorf("invalid password hash format")
	}

	if parts[1] != "argon2id" {
		return false, fmt.Errorf("unsupported password hash algorithm")
	}

	// Parse Argon2 parameters.
	parameters := strings.Split(parts[3], ",")

	if len(parameters) != 3 {
		return false, fmt.Errorf("invalid Argon2 parameters")
	}

	var memory uint32
	var time uint32
	var threads uint8

	for _, parameter := range parameters {
		value := strings.Split(parameter, "=")

		if len(value) != 2 {
			return false, fmt.Errorf("invalid Argon2 parameter")
		}

		number, err := strconv.ParseUint(value[1], 10, 32)
		if err != nil {
			return false, fmt.Errorf("invalid Argon2 parameter value")
		}

		switch value[0] {
		case "m":
			memory = uint32(number)

		case "t":
			time = uint32(number)

		case "p":
			threads = uint8(number)

		default:
			return false, fmt.Errorf("unknown Argon2 parameter")
		}
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode password hash: %w", err)
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(actualHash, expectedHash) == 1 {
		return true, nil
	}

	return false, nil
}