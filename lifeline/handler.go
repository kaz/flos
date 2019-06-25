package lifeline

import "github.com/labstack/echo/v4"

func get(c echo.Context) error {
	c.Set("response", results)
	return nil
}
