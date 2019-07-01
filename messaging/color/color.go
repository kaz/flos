package color

import (
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/kaz/flos/camo"
	"github.com/shamaton/msgpack"
)

const (
	// signature valid in 15s
	VALID_THRU = 15 * time.Second

	SIGN_KEY = "Daphne Ficus Iris Maackia Lythrum Myrica Sabia Flos"
)

type (
	Protocol struct{}

	StampedPayload struct {
		Payload   []byte
		Timestamp int64
	}
	SignedPayload struct {
		Payload   []byte
		Signature []byte
	}
)

func serialize(obj interface{}) ([]byte, error) {
	return msgpack.Encode(obj)
}
func deserialize(data []byte, objPtr interface{}) error {
	return msgpack.Decode(data, objPtr)
}

func sign(data []byte) ([]byte, error) {
	stamped, err := serialize(&StampedPayload{
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

	return serialize(&SignedPayload{
		stamped,
		m.Sum(nil),
	})
}
func verify(data []byte) ([]byte, error) {
	signed := &SignedPayload{}
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

	stamped := &StampedPayload{}
	if err := deserialize(signed.Payload, stamped); err != nil {
		return nil, err
	}
	if time.Since(time.Unix(0, stamped.Timestamp)) > VALID_THRU {
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
