package archive

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"go.etcd.io/bbolt"
)

const (
	DB_FILE = "chunk.0003.zip"
)

var (
	db *bbolt.DB

	logger = log.New(os.Stdout, "[archive] ", log.Ltime)
)

func StartService(g *echo.Group) {
	var err error
	db, err = bbolt.Open(DB_FILE, 0644, nil)
	if err != nil {
		logger.Printf("Failed to open db: %v\n", err)
		return
	}

	go runMaster()
}
