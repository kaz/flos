package archive

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

const (
	DB_FILE = "chunk.0003.zip"
)

var (
	logger = log.New(os.Stdout, "[archive] ", log.Ltime)
)

func StartService(g *echo.Group) {
	go runMaster()
}
