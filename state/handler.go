package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo/v4"
)

func getConfig(c echo.Context) error {
	mu.RLock()
	defer mu.RUnlock()

	c.Set("response", store)
	return nil
}

func putConfig(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()

	req, ok := c.Get("request").(map[string]interface{})
	if !ok {
		return fmt.Errorf("unpected request format")
	}

	raw, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(STORE_FILE, raw, 0644); err != nil {
		return err
	}

	store = req
	rawStore = raw

	return nil
}
