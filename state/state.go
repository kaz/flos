package state

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/labstack/echo/v4"
)

const (
	STORE_FILE = "./state.json"
)

var (
	mu = sync.RWMutex{}

	store    interface{}
	rawStore []byte
)

func StartService(g *echo.Group) {
	var err error
	rawStore, err = ioutil.ReadFile(STORE_FILE)
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
