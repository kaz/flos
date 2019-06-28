package color

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/kaz/flos/camo"
)

const (
	// signature valid in 4s
	VALID_THRU = 4 * 1000 * 1000

	SIGN_KEY = "Daphne Ficus Iris Maackia Lythrum Myrica Sabia Flos"
)

type (
	Protocol struct{}

	stampedPayload struct {
		Payload   []byte
		Timestamp int64
	}
	signedPayload struct {
		Payload   []byte
		Signature []byte
	}
)

func serialize(obj interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(obj); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func deserialize(data []byte, objPtr interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(objPtr)
}

func sign(data []byte) ([]byte, error) {
	stamped, err := serialize(stampedPayload{
		data,
		time.Now().UnixNano(),
	})
	if err != nil {
		return nil, err
	}

	m := hmac.New(md5.New, []byte(SIGN_KEY))
	if _, err := m.Write(stamped); err != nil {
		return nil, err
	}

	return serialize(signedPayload{
		stamped,
		m.Sum(nil),
	})
}
func verify(data []byte) ([]byte, error) {
	signed := &signedPayload{}
	if err := deserialize(data, signed); err != nil {
		return nil, err
	}

	m := hmac.New(md5.New, []byte(SIGN_KEY))
	if _, err := m.Write(signed.Payload); err != nil {
		return nil, err
	}

	if !hmac.Equal(signed.Signature, m.Sum(nil)) {
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

func (p *Protocol) Encode(obj interface{}) ([]byte, error) {
	data, err := serialize(obj)
	if err != nil {
		return nil, err
	}
	data, err = sign(data)
	if err != nil {
		return nil, err
	}
	return camo.Encode(data)
}

func (p *Protocol) Decode(data []byte, objPtr interface{}) error {
	data, err := camo.Decode(data)
	if err != nil {
		return err
	}
	data, err = verify(data)
	if err != nil {
		return err
	}
	return deserialize(data, objPtr)
}

func (p *Protocol) Type() string {
	return "application/octet-stream"
}
