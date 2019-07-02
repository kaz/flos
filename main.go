package main

import (
	"github.com/kaz/flos/archive"
	"github.com/kaz/flos/audit"
	"github.com/kaz/flos/libra"
	"github.com/kaz/flos/lifeline"
	"github.com/kaz/flos/limit"
	"github.com/kaz/flos/messaging"
	"github.com/kaz/flos/power"
	"github.com/kaz/flos/proxy"
	"github.com/kaz/flos/state"
	"github.com/kaz/flos/tail"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = messaging.ErrorHandler

	limit.Init()
	messaging.Init()

	e.Pre(proxy.Middleware)
	e.Pre(messaging.Middleware)

	power.StartService(e.Group("/power"))
	libra.StartService(e.Group("/libra"))
	state.StartService(e.Group("/state"))
	archive.StartService(e.Group("/archive"))
	lifeline.StartService(e.Group("/lifeline"))

	go tail.StartWorker()
	go audit.StartWorker()

	e.Logger.Fatal(e.Start(power.LISTEN))
}
