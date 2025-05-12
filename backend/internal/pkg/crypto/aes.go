package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	pbkdf2Iterations = 4096
	saltSize         = 16
	keySize          = 32
)

type aesCipher struct {
	key []byte
}

func NewAESCipher(key []byte) Cipher {
	return &aesCipher{key: key}
}

func generateRandomSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	return salt, nil
}

func (p *aesCipher) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", ErrCipherTextEmpty
	}

	salt, err := generateRandomSalt(saltSize)
	if err != nil {
		return "", err
	}

	key := pbkdf2.Key(p.key, salt, pbkdf2Iterations, keySize, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	ciphertextBytes := append(salt, encrypted...)
	return base64.StdEncoding.EncodeToString(ciphertextBytes), nil
}

func (p *aesCipher) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", ErrCipherTextEmpty
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < saltSize {
		return "", ErrCipherTextTooShort
	}

	salt := ciphertextBytes[:saltSize]
	encrypted := ciphertextBytes[saltSize:]
	key := pbkdf2.Key(p.key, salt, pbkdf2Iterations, keySize, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return "", ErrCipherTextTooShort
	}

	nonce, ct := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
