package libra

import (
	"fmt"
	"log"
	"os"

	"github.com/kaz/flos/libra/bookshelf"
	"github.com/labstack/echo/v4"
)

const (
	LIBRA_FILE = "chunk.0002.zip"
)

var (
	logger = log.New(os.Stdout, "[libra] ", log.Ltime)
	libra  *bookshelf.Bookshelf
)

func StartService(g *echo.Group) {
	lib, err := bookshelf.New(LIBRA_FILE)
	if err != nil {
		logger.Printf("Failed to open db: %v\n", err)
		return
	}

	libra = lib

	g.GET("", libra.ListHandler)
	g.PATCH("/books", libra.GetHandler)
	g.DELETE("/books", libra.DeleteHandler)
}

func Put(series, contents string) error {
	if libra == nil {
		return fmt.Errorf("library is closed")
	}
	return libra.Put([]byte(series), []byte(contents))
}
