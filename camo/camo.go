package camo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
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

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	aead, err := gcm()
	if err != nil {
		return err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	contents := []byte(CAMO_HEADER)
	contents = append(contents, nonce...)
	contents = append(contents, aead.Seal(nil, nonce, data, nil)...)

	return ioutil.WriteFile(filename, contents, perm)
}
func ReadFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	aead, err := gcm()
	if err != nil {
		return nil, err
	}

	size := aead.NonceSize()
	if size+len(CAMO_HEADER) > len(data) {
		return nil, fmt.Errorf("file size too short")
	}
	data = data[len(CAMO_HEADER):]

	return aead.Open(nil, data[:size], data[size:], nil)
}
