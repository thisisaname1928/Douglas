package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

// encryption and decryption work by using AES-128 + GCM

func Encrypt(data []byte, key string) ([]byte, error) {
	if len(key) > 16 {
		return nil, errors.New("key too long")
	}
	for len(key) < 16 {
		key += "-"
	}
	block, e := aes.NewCipher([]byte(key))
	if e != nil {
		return []byte{}, e
	}

	gcm, e := cipher.NewGCM(block)
	if e != nil {
		return []byte{}, e
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, e := io.ReadFull(rand.Reader, nonce); e != nil {
		return []byte{}, e
	}

	encryptedContent := gcm.Seal(nonce, nonce, data, nil)

	return encryptedContent, nil
}

func Decrypt(data []byte, key string) ([]byte, error) {
	if len(key) > 16 {
		return nil, errors.New("key too long")
	}
	for len(key) < 16 {
		key += "-"
	}

	block, e := aes.NewCipher([]byte(key))
	if e != nil {
		return []byte{}, e
	}

	gcm, e := cipher.NewGCM(block)
	if e != nil {
		return []byte{}, e
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return []byte{}, errors.New("can't decrypt data")
	}
	nonce := data[:nonceSize]
	encryptedData := data[nonceSize:]

	decryptedData, e := gcm.Open(nil, nonce, encryptedData, nil)

	return decryptedData, e
}

// should use bcrypt instead
func EncryptKey(key string) string {
	hsh := sha256.New()
	hsh.Write([]byte(key))

	outputKey := hsh.Sum(nil)

	return fmt.Sprintf("%x", outputKey)
}
