package messaging

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	logger.Println(err)
	resp, err := Encode(err.Error())
	if err != nil {
		logger.Println("failed to encode:", err)
	}
	c.Blob(http.StatusBadRequest, Type(), resp)
}

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		c.Request().Body.Close()

		var req interface{}
		if len(body) > 0 {
			if err := Decode(body, &req); err != nil {
				return err
			}
		}
		c.Set("request", req)

		if err := next(c); err != nil {
			return err
		}

		resp, err := Encode(c.Get("response"))
		if err != nil {
			return err
		}

		return c.Blob(http.StatusOK, Type(), resp)
	}
}
