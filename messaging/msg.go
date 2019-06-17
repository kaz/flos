package messaging

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

const (
	KEY_SIGN = "Daphne Ficus Iris Maackia"
	KEY_ENC  = "Lythrum Myrica Sabia Flos"

	// signature valid in 5s
	VALID_THRU = 5 * 1000 * 1000
)

type (
	stampedPayload struct {
		Payload   interface{}
		Timestamp int64
	}
	signedPayload struct {
		Payload   []byte
		Signature []byte
	}
	encryptedPayload struct {
		Payload []byte
		Nonce   []byte
	}
)

func serialize(obj interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 256*1024))
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func deserialize(data []byte, objPtr interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(objPtr)
}

func sign(obj interface{}) ([]byte, error) {
	stamped, err := serialize(stampedPayload{
		obj,
		time.Now().UnixNano(),
	})
	if err != nil {
		return nil, err
	}

	return serialize(signedPayload{
		stamped,
		hmac.New(sha256.New, []byte(KEY_SIGN)).Sum(stamped),
	})
}
func verify(data []byte) (interface{}, error) {
	signed := &signedPayload{}
	if err := deserialize(data, signed); err != nil {
		return nil, err
	}
	if !hmac.Equal(signed.Signature, hmac.New(sha256.New, []byte(KEY_SIGN)).Sum(signed.Payload)) {
		return nil, fmt.Errorf("signature not match")
	}

	stamped := &stampedPayload{}
	if err := deserialize(signed.Payload, stamped); err != nil {
		return nil, err
	}
	if time.Now().UnixNano()-stamped.Timestamp > VALID_THRU {
		return nil, fmt.Errorf("signature expired")
	}

	return stamped.Payload, nil
}

func encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(KEY_ENC[:aes.BlockSize]))
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return serialize(encryptedPayload{
		aead.Seal(nil, nonce, data, nil),
		nonce,
	})
}
func decrypt(data []byte) ([]byte, error) {
	encrypted := &encryptedPayload{}
	if err := deserialize(data, encrypted); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(KEY_ENC[:aes.BlockSize]))
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, encrypted.Nonce, encrypted.Payload, nil)
}

func Encode(obj interface{}) ([]byte, error) {
	data, err := sign(obj)
	if err != nil {
		return nil, err
	}
	return encrypt(data)
}

func Decode(data []byte) (interface{}, error) {
	data, err := decrypt(data)
	if err != nil {
		return nil, err
	}
	return verify(data)
}
