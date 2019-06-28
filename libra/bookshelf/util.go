package bookshelf

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"encoding/gob"
	"sync"
)

var (
	mu     sync.RWMutex
	shelfs = []string{}
)

func registerBookshelf(path string) {
	mu.Lock()
	defer mu.Unlock()

	shelfs = append(shelfs, path)
}
func IsBookshelf(path string) bool {
	mu.RLock()
	defer mu.RUnlock()

	for _, s := range shelfs {
		if path == s {
			return true
		}
	}
	return false
}

func itob(i uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, i)
	return key
}
func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

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

func hash(key []byte) ([]byte, error) {
	m := hmac.New(md5.New, []byte(HASH_KEY))
	if _, err := m.Write(key); err != nil {
		return nil, err
	}
	return m.Sum(nil), nil
}
