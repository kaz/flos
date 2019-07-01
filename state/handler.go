package state

import (
	"github.com/shamaton/msgpack"

	"github.com/kaz/flos/camo"
	"github.com/labstack/echo/v4"
)

func getConfig(c echo.Context) error {
	mu.RLock()
	c.Set("response", current)
	mu.RUnlock()

	return nil
}

func putConfig(c echo.Context) error {
	raw, err := msgpack.Encode(c.Get("request"))
	if err != nil {
		return err
	}

	var newState State
	if err := msgpack.Decode(raw, &newState); err != nil {
		return err
	}

	if err := camo.WriteFile(STORE_FILE, raw, 0644); err != nil {
		return err
	}

	mu.Lock()
	current = newState
	mu.Unlock()

	logger.Println("state updated")
	return nil
}
