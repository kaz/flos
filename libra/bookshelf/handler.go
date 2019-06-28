package bookshelf

import (
	"bytes"
	"fmt"

	"github.com/kaz/flos/camo"
	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
)

func (b *Bookshelf) ListHandler(c echo.Context) error {
	data, err := b.listSeries()
	if err != nil {
		return err
	}

	c.Set("response", data)
	return nil
}
func (b *Bookshelf) GetHandler(c echo.Context) error {
	req, ok := c.Get("request").(float64)
	if !ok {
		return fmt.Errorf("unexpected request format")
	}

	data, err := b.getAfter(uint64(req))
	if err != nil {
		return err
	}

	c.Set("response", data)
	return nil
}
func (b *Bookshelf) DeleteHandler(c echo.Context) error {
	req, ok := c.Get("request").(float64)
	if !ok {
		return fmt.Errorf("unexpected request format")
	}

	if err := b.deleteBefore(uint64(req)); err != nil {
		return err
	}

	return nil
}

func (b *Bookshelf) listSeries() ([][]byte, error) {
	result := [][]byte{}

	return result, b.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket([]byte(META_BUCKET)).Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			data, err := camo.Decode(k)
			if err != nil {
				return err
			}

			result = append(result, data)
		}

		return nil
	})
}
func (b *Bookshelf) getAfter(gte uint64) ([]*Book, error) {
	result := []*Book{}

	return result, b.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket([]byte(DATA_BUCKET)).Cursor()

		for k, v := cursor.Seek(itob(gte)); k != nil; k, v = cursor.Next() {
			data, err := camo.Decode(v)
			if err != nil {
				return err
			}

			var book *Book
			if err := deserialize(data, book); err != nil {
				return err
			}

			book.ID = btoi(k)
			result = append(result, book)

			if len(result) >= MAX_ROW_COUNT {
				break
			}
		}

		return nil
	})
}
func (b *Bookshelf) deleteBefore(lte uint64) error {
	end := itob(lte)

	return b.db.Update(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket([]byte(DATA_BUCKET)).Cursor()

		for k, _ := cursor.First(); k != nil && bytes.Compare(k, end) < 1; k, _ = cursor.Next() {
			if err := cursor.Delete(); err != nil {
				return err
			}
		}

		return nil
	})
}
