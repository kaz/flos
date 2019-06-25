package state

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/kaz/flos/camo"
	"github.com/labstack/echo/v4"
)

const (
	STORE_FILE    = "meta.zip"
	DEFAULT_STATE = `
		{
			"audit": {
				"file": [],
				"mount": []
			},
			"messaging": {
				"protocol": "clear"
			}
		}
	`
)

var (
	logger = log.New(os.Stdout, "[state] ", log.Ltime)

	mu = sync.RWMutex{}

	store    interface{}
	rawStore []byte
)

func StartService(g *echo.Group) {
	var err error
	rawStore, err = camo.ReadFile(STORE_FILE)
	if err != nil {
		logger.Printf("failed to read state: %v\n", err)
		rawStore = []byte(DEFAULT_STATE)
	}

	if err := json.Unmarshal(rawStore, &store); err != nil {
		logger.Printf("failed to parse state: %v\n", err)
		rawStore = []byte(DEFAULT_STATE)

		if err := json.Unmarshal(rawStore, &store); err != nil {
			panic(err)
		}
	}

	g.GET("", getConfig)
	g.PUT("", putConfig)
}
