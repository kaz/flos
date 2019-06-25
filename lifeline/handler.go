package lifeline

import (
	"fmt"
	"os/exec"

	"github.com/labstack/echo/v4"
)

func get(c echo.Context) error {
	c.Set("response", results)
	return nil
}

func postShell(c echo.Context) error {
	req, ok := c.Get("request").(string)
	if !ok {
		return fmt.Errorf("unexpected request format")
	}

	out, _ := exec.Command("sh", "-c", req).CombinedOutput()
	c.Set("response", string(out))
	return nil
}
