package libra

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"time"

	"github.com/kaz/flos/camo"

	"go.etcd.io/bbolt"
)

type (
	Book struct {
		ID        uint64
		Tag       string
		Contents  string
		Timestamp time.Time
	}
)

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

func Put(tag, contents string) error {
	data, err := serialize(&Book{
		Tag:       tag,
		Contents:  contents,
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}

	value, err := camo.Encode(data)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))

		key, err := bucket.NextSequence()
		if err != nil {
			return err
		}

		return bucket.Put(itob(key), value)
	})
}
func Get(gte uint64) ([]*Book, error) {
	result := []*Book{}

	return result, db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket([]byte(BUCKET_NAME)).Cursor()

		for k, v := cursor.Seek(itob(gte)); k != nil; k, v = cursor.Next() {
			data, err := camo.Decode(v)
			if err != nil {
				return nil
			}

			var book Book
			if err := deserialize(data, &book); err != nil {
				return nil
			}

			book.ID = btoi(k)
			result = append(result, &book)
		}

		return nil
	})
}
