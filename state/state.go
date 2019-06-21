package state

import (
	"encoding/json"
	"sync"

	"github.com/kaz/flos/camo"
	"github.com/labstack/echo/v4"
)

const (
	STORE_FILE = "meta.zip"
)

var (
	mu = sync.RWMutex{}

	store    interface{}
	rawStore []byte
)

func StartService(g *echo.Group) {
	var err error
	rawStore, err = camo.ReadFile(STORE_FILE)
	if err != nil {
		rawStore = []byte(`{"_state":"created"}`)
	}

	if err := json.Unmarshal(rawStore, &store); err != nil {
		rawStore = []byte(`{"_state":"discarded"}`)
		if err := json.Unmarshal(rawStore, &store); err != nil {
			panic(err)
		}
	}

	g.GET("", getConfig)
	g.PUT("", putConfig)
}
