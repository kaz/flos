package beacon

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/kaz/flos/power"
	"github.com/labstack/echo/v4"
)

const (
	BEACON_CYCLE_SEC = 5
	UDP_ADDR         = "239.239.239.239" + power.LISTEN
)

var (
	logger = log.New(os.Stdout, "[beacon] ", log.Ltime)

	mu    = sync.RWMutex{}
	nodes = map[string]time.Time{}
)

func StartService(g *echo.Group) {
	go startSender()
	go startReceiver()

	g.GET("/nodes", getNodes)
	g.DELETE("/node", deleteNode)
}
