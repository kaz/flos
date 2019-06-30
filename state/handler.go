package state

import (
	"encoding/json"
	"fmt"

	"github.com/kaz/flos/camo"
	"github.com/labstack/echo/v4"
)

func getConfig(c echo.Context) error {
	mu.RLock()
	defer mu.RUnlock()

	resp := make(map[string]interface{}, len(store))
	for k, v := range store {
		resp[k] = v
	}

	c.Set("response", resp)
	return nil
}

func putConfig(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()

	req, ok := c.Get("request").(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected request format")
	}

	raw, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	if err := camo.WriteFile(STORE_FILE, raw, 0644); err != nil {
		return err
	}

	store = req
	rawStore = raw

	logger.Println("state updated")
	return nil
}
