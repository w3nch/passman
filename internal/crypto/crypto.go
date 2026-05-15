package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	KeySize    = 32
	SaltSize   = 32
	Iterations = 100000
)

func GenerateSalt() string {
	salt := make([]byte, SaltSize)
	rand.Read(salt)
	return base64.StdEncoding.EncodeToString(salt)
}

func DeriveKey(password, salt string) []byte {
	saltBytes, _ := base64.StdEncoding.DecodeString(salt)
	return pbkdf2.Key([]byte(password), saltBytes, Iterations, KeySize, sha256.New)
}

func HashPassword(password, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func GenerateRecoveryCodes(count int) []string {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		code := make([]byte, 16)
		rand.Read(code)
		codes[i] = base64.URLEncoding.EncodeToString(code)[:16]
	}
	return codes
}

func VerifyPassword(password, salt, hash string) bool {
	return HashPassword(password, salt) == hash
}

func MaskPassword(pwd string) string {
	return fmt.Sprintf("%s%s", pwd[:2], "********")
}

func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func Decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}