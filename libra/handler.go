package libra

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func getBooksAfter(c echo.Context) error {
	req, ok := c.Get("request").(float64)
	if !ok {
		return fmt.Errorf("unexpected request format")
	}

	data, err := getAfter(uint64(req))
	if err != nil {
		return err
	}

	c.Set("response", data)
	return nil
}
func deleteBooksBefore(c echo.Context) error {
	req, ok := c.Get("request").(float64)
	if !ok {
		return fmt.Errorf("unexpected request format")
	}

	if err := deleteBefore(uint64(req)); err != nil {
		return err
	}

	return nil
}
