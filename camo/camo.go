package camo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/DataDog/zstd"
)

const (
	CAMO_KEY    = "Thymus Ribes Abelia Sedum Felicia Ochna Lychnis"
	CAMO_HEADER = "\x50\x4b\x03\x04\x0a\x00\x00\x00\x00\x00\xc0\xb4\xd4\x4e\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x1c\x00\x5f\x55\x54\x09\x00\x03\x37\x8c"
)

func gcm() (cipher.AEAD, error) {
	block, err := aes.NewCipher([]byte(CAMO_KEY[:32]))
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(block)
}
func encrypt(data []byte) ([]byte, error) {
	aead, err := gcm()
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return append(nonce, aead.Seal(nil, nonce, data, nil)...), nil
}
func decrypt(data []byte) ([]byte, error) {
	aead, err := gcm()
	if err != nil {
		return nil, err
	}

	size := aead.NonceSize()
	if size > len(data) {
		return nil, fmt.Errorf("invalid payload length")
	}

	return aead.Open(nil, data[:size], data[size:], nil)
}

func compress(data []byte) ([]byte, error) {
	return zstd.Compress(nil, data)
}
func decompress(data []byte) ([]byte, error) {
	return zstd.Decompress(nil, data)
}
