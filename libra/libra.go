package libra

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
)

const (
	DB_FILE     = "chunk.zip"
	BUCKET_NAME = "PK"
)

var (
	db *bbolt.DB

	logger = log.New(os.Stdout, "[libra] ", log.Ltime)
)

func StartService(g *echo.Group) {
	var err error
	db, err = bbolt.Open(DB_FILE, 0644, nil)
	if err != nil {
		logger.Printf("Failed to open db: %v\n", err)
		return
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))
		return err
	})
	if err != nil {
		logger.Printf("Failed to create bucket: %v\n", err)
		return
	}

	g.PATCH("/books", getBooksAfter)
	g.DELETE("/books", deleteBooksBefore)
}
