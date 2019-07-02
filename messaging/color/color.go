package color

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/kaz/flos/camo"
	"github.com/shamaton/msgpack"
)

const (
	// signature valid in 15s
	VALID_THRU = 15 * time.Second

	SIGN_KEY = "Lythrum Myrica Sabia Flos"
)

type (
	Protocol struct{}
)

func itob(i int64) []byte {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(i))
	return key
}
func btoi(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}

func serialize(obj interface{}) ([]byte, error) {
	return msgpack.Encode(obj)
}
func deserialize(data []byte, objPtr interface{}) error {
	return msgpack.Decode(data, objPtr)
}

func sign(data []byte) ([]byte, error) {
	data = append(data, itob(time.Now().UnixNano())...)

	m := hmac.New(md5.New, []byte(SIGN_KEY))
	if _, err := m.Write(data); err != nil {
		return nil, err
	}

	return append(data, m.Sum(nil)...), nil
}
func verify(data []byte) ([]byte, error) {
	p := len(data) - md5.Size
	if p <= 0 {
		return nil, fmt.Errorf("invalid size")
	}

	signature := data[p:]
	data = data[:p]

	m := hmac.New(md5.New, []byte(SIGN_KEY))
	if _, err := m.Write(data); err != nil {
		return nil, err
	}

	if !hmac.Equal(m.Sum(nil), signature) {
		return nil, fmt.Errorf("signature not match")
	}

	p = len(data) - 8
	if p <= 0 {
		return nil, fmt.Errorf("invalid size")
	}

	ts := data[p:]
	data = data[:p]

	if time.Since(time.Unix(0, btoi(ts))) > VALID_THRU {
		return nil, fmt.Errorf("signature expired")
	}

	return data, nil
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
