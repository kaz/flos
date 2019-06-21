package libra

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func getBooksAfter(c echo.Context) error {
	req, ok := c.Get("request").(float64)
	if !ok {
		return fmt.Errorf("unpected request format")
	}

	data, err := GetAfter(uint64(req))
	if err != nil {
		return err
	}

	c.Set("response", data)
	return nil
}
func deleteBooksBefore(c echo.Context) error {
	req, ok := c.Get("request").(float64)
	if !ok {
		return fmt.Errorf("unpected request format")
	}

	if err := DeleteBefore(uint64(req)); err != nil {
		return err
	}

	return nil
}
