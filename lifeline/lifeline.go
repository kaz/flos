package lifeline

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

var (
	logger = log.New(os.Stdout, "[lifeline] ", log.Ltime)

	status = map[string]*result{}
)

func StartService(g *echo.Group) {
	go runMaster()

	g.GET("", get)
}
