package power

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func postPower(c echo.Context) error {
	req, ok := c.Get("request").(string)
	if !ok {
		return fmt.Errorf("unexpected request format")
	}

	delay := func(action func()) {
		time.Sleep(DELAY_SEC * time.Second)
		action()
	}

	if req == "stop" {
		go delay(stop)
	} else if req == "restart" {
		go delay(restart)
	} else {
		return fmt.Errorf("no such action: %v\n", req)
	}

	return nil
}
