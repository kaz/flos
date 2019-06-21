package libra

import (
	"log"
	"os"

	"go.etcd.io/bbolt"
)

const (
	DB_FILE     = "chunk.zip"
	BUCKET_NAME = ""
)

var (
	db *bbolt.DB

	logger = log.New(os.Stdout, "[libra] ", log.Ltime)
)

func StartService() {
	var err error
	db, err = bbolt.Open(DB_FILE, 0644, nil)
	if err != nil {
		logger.Printf("Failed to open db: %v\n", err)
		return
	}
}
