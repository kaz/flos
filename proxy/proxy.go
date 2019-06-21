package proxy

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kaz/flos/messaging"

	"github.com/labstack/echo/v4"
)

var (
	logger = log.New(os.Stdout, "[proxy] ", log.Ltime)
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			return next(c)
		}

		spld := strings.Split(authHeader, " ")
		if len(spld) != 2 || spld[0] != "bearer" {
			return next(c)
		}

		payload, err := base64.StdEncoding.DecodeString(spld[1])
		if err != nil {
			logger.Printf("Request rejected: %v\n", err)
			return next(c)
		}

		var destination string
		if err := messaging.Decode(payload, &destination); err != nil {
			logger.Printf("Request rejected: %v\n", err)
			return next(c)
		}

		proxyReq, err := http.NewRequest(req.Method, c.Scheme()+"://"+destination+req.URL.Path, req.Body)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			return err
		}

		return c.Stream(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Body)
	}
}
