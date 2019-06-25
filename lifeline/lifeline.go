package lifeline

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

var (
	logger = log.New(os.Stdout, "[lifeline] ", log.Ltime)

	results = map[string]*Result{}
)

func StartService(g *echo.Group) {
	go runMaster()

	g.GET("", get)
}
