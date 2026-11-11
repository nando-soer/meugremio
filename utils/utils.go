package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

func randomizeBytes(length uint32) (byteSlice []byte, err error) {
	byteSlice = make([]byte, length)
	_, err = rand.Read(byteSlice)

	if err != nil {
		return nil, err
	}

	return
}

func GenerateSessionId() (string, error) {
	random, err := randomizeBytes(32)

	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(random), err
}

func GenerateConfirmationCode() (string, error) {
	random, err := randomizeBytes(5)

	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(random), err
}

func GenerateRecoveryToken() (string, error) {
	return GenerateSessionId() // Mesmo algoritmo
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	storableHash := hex.EncodeToString(hash[:])
	return storableHash
}

func CompareTokens(receivedToken, encodedToken string) bool {
	receivedHash := HashToken(receivedToken)
	return subtle.ConstantTimeCompare([]byte(receivedHash), []byte(encodedToken)) == 1
}

func encryptPassword(password, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	defer h.Reset()
	h.Write([]byte(password))
	return h.Sum(nil)
}

func generateHash(password string, salt []byte, p *params) []byte {
	pepper := os.Getenv("PEPPER")
	pepperedPassword := encryptPassword(password, pepper)

	return argon2.IDKey(pepperedPassword, salt, p.iterations, p.memory, p.parallelism, p.keyLength)
}

func encodeHash(hash, salt []byte, p *params) string {
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, encodedSalt, encodedHash)
}

func decodeHash(encodedString string) (hash []byte, salt []byte, p params, err error) {
	segments := strings.Split(encodedString, "$")

	if segments[1] != "argon2id" {
		return nil, nil, params{}, errors.New("Incompatible hashing algorithm")
	}

	var version int
	_, err = fmt.Sscanf(segments[2], "v=%d", &version)
	if err != nil {
		return nil, nil, params{}, err
	}

	if version > argon2.Version {
		return nil, nil, params{}, errors.New("Incompatible Argon2 version")
	}

	_, err = fmt.Sscanf(segments[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, params{}, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(segments[4])
	if err != nil {
		return nil, nil, params{}, err
	}

	hash, err = base64.RawStdEncoding.Strict().DecodeString(segments[5])
	if err != nil {
		return nil, nil, params{}, err
	}

	p.keyLength = uint32(len(hash))

	return
}

func HashPassword(password string) (encodedString string, err error) {
	// REFERENCE: https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
	defaultParams := params{
		memory:      64 * 1024,
		iterations:  2,
		parallelism: 1,
		keyLength:   32,
	}

	const saltLength = 16
	salt, err := randomizeBytes(saltLength)
	if err != nil {
		return "", err
	}

	hash := generateHash(password, salt, &defaultParams)
	encodedString = encodeHash(hash, salt, &defaultParams)

	return
}

func ComparePasswords(input, encodedPassword string) (bool, error) {
	passwordHash, passwordSalt, params, err := decodeHash(encodedPassword)
	if err != nil {
		return false, err
	}

	inputHash := generateHash(input, passwordSalt, &params)

	return subtle.ConstantTimeCompare(inputHash, passwordHash) == 1, nil
}
