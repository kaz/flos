package camo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pierrec/lz4"
)

const (
	CAMO_KEY    = "Daphne Ficus Iris Maackia"
	CAMO_HEADER = "\x50\x4b\x03\x04\x0a\x00\x00\x00\x00\x00\xc0\xb4\xd4\x4e\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x1c\x00\x5f\x55\x54\x09\x00\x03\x37\x8c"
)

func getStream(iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher([]byte(CAMO_KEY[:aes.BlockSize]))
	if err != nil {
		return nil, err
	}

	return cipher.NewCTR(block, iv), nil
}
func encrypt(w io.Writer) (io.Writer, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	stream, err := getStream(iv)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(iv); err != nil {
		return nil, err
	}

	return &cipher.StreamWriter{S: stream, W: w}, nil
}
func decrypt(r io.Reader) (io.Reader, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := r.Read(iv); err != nil {
		return nil, err
	}

	stream, err := getStream(iv)
	if err != nil {
		return nil, err
	}

	return &cipher.StreamReader{S: stream, R: r}, nil
}

func compress(w io.Writer) *lz4.Writer {
	return lz4.NewWriter(w)
}
func decompress(r io.Reader) *lz4.Reader {
	return lz4.NewReader(r)
}
