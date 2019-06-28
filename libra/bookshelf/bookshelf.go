package bookshelf

import (
	"path/filepath"
	"time"

	"github.com/kaz/flos/camo"
	"go.etcd.io/bbolt"
)

const (
	META_BUCKET = "META"
	DATA_BUCKET = "DATA"
)

type (
	Bookshelf struct {
		DBFile string
		db     *bbolt.DB
		maxRow int
	}
	Book struct {
		ID        uint64
		Series    []byte
		Contents  []byte
		Timestamp time.Time
	}
)

func New(path string, maxRow int) (*Bookshelf, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	db, err := bbolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(META_BUCKET)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(DATA_BUCKET)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	registerBookshelf(abspath)
	return &Bookshelf{abspath, db, maxRow}, nil
}

func (b *Bookshelf) Put(series, contents []byte) error {
	data, err := serialize(&Book{
		Series:    series,
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

	encodedSeries, err := camo.Encode(series)
	if err != nil {
		return err
	}

	return b.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.Bucket([]byte(META_BUCKET)).Put(encodedSeries, []byte{}); err != nil {
			return err
		}

		bucket := tx.Bucket([]byte(DATA_BUCKET))

		key, err := bucket.NextSequence()
		if err != nil {
			return err
		}

		return bucket.Put(itob(key), value)
	})
}

func (b *Bookshelf) Has(series []byte) (bool, error) {
	encodedSeries, err := camo.Encode(series)
	if err != nil {
		return false, err
	}

	var has bool
	return has, b.db.View(func(tx *bbolt.Tx) error {
		has = tx.Bucket([]byte(META_BUCKET)).Get(encodedSeries) != nil
		return nil
	})
}
