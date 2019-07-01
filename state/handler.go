package state

import (
	"fmt"

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
	req, ok := c.Get("request").(State)
	if !ok {
		return fmt.Errorf("unexpected request format")
	}

	raw, err := msgpack.Encode(req)
	if err != nil {
		return err
	}

	if err := camo.WriteFile(STORE_FILE, raw, 0644); err != nil {
		return err
	}

	mu.Lock()
	current = req
	mu.Unlock()

	logger.Println("state updated")
	return nil
}
