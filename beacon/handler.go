package beacon

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func getNodes(c echo.Context) error {
	mu.RLock()
	defer mu.RUnlock()

	resp := make(map[string]time.Time)
	for k, v := range nodes {
		resp[k] = v
	}

	c.Set("response", resp)
	return nil
}
func deleteNode(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()

	req, ok := c.Get("request").(string)
	if !ok {
		return fmt.Errorf("unexpected request format")
	}

	delete(nodes, req)
	return nil
}
