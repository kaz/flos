package main

import (
	"log"

	"github.com/kaz/flos/beacon"
	"github.com/kaz/flos/messaging"

	"github.com/labstack/echo/v4"
)

func logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			log.Printf("[main] internal error: %v\n", err)
		}
		return err
	}
}

func main() {
	e := echo.New()
	e.Use(logger)
	e.Use(messaging.Middleware)

	beacon.StartService(e.Group("/beacon"))

	e.Logger.Fatal(e.Start(":39239"))
}
