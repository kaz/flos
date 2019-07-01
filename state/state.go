package state

import (
	"log"
	"os"
	"sync"

	"github.com/kaz/flos/camo"
	"github.com/labstack/echo/v4"
	"github.com/shamaton/msgpack"
)

const (
	STORE_FILE = "chunk.0001.zip"
)

var (
	logger = log.New(os.Stdout, "[state] ", log.Ltime)
	mu     = sync.RWMutex{}
)

func StartService(g *echo.Group) {
	raw, err := camo.ReadFile(STORE_FILE)
	if err != nil {
		logger.Printf("failed to read state: %v\n", err)
	} else if err := msgpack.Decode(raw, &current); err != nil {
		logger.Printf("failed to parse state: %v\n", err)
	}

	g.GET("", getConfig)
	g.PUT("", putConfig)
}

func Get() State {
	return current
}
