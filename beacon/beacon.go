package beacon

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	BEACON_CYCLE_SEC = 6
	UDP_ADDR         = "239.239.239.239:239"
)

var (
	logger = log.New(os.Stdout, "[beacon] ", log.Ltime)

	mu    = sync.RWMutex{}
	nodes = map[string]time.Time{}
)

func StartService(g *echo.Group) {
	go func() {
		for {
			ch := make(chan error)
			go sendBeacon(ch)
			logger.Printf("Sending beacon failed: %v\n", <-ch)
			close(ch)
		}
	}()
	go func() {
		for {
			ch := make(chan error)
			go recvBeacon(ch)
			logger.Printf("Receiving beacon failed: %v\n", <-ch)
			close(ch)
		}
	}()

	g.GET("/nodes", getNodes)
	g.DELETE("/node", deleteNode)
}
