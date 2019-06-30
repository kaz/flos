package lifeline

import (
	"log"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	logger = log.New(os.Stdout, "[lifeline] ", log.Ltime)

	mu      = sync.RWMutex{}
	results = map[string]*Result{}
)

func StartService(g *echo.Group) {
	go runMaster()

	g.GET("", get)
	g.POST("/shell", postShell)
}
