package power

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/kaz/flos/messaging"
)

const (
	ACTION_DELAY_SEC = 3
	LISTEN           = ":39239"
)

var (
	logger = log.New(os.Stdout, "[power] ", log.Ltime)
)

func StartService(g *echo.Group) {
	for {
		lis, err := net.Listen("tcp", LISTEN)
		if err == nil {
			lis.Close()
			break
		}

		logger.Println("Trying to kill other process ...")

		payload, err := messaging.Encode("stop")
		if err != nil {
			logger.Printf("Failed to encode payload: %v\n", err)
			time.Sleep(ACTION_DELAY_SEC * time.Second)
			continue
		}

		resp, err := http.Post("http://"+LISTEN+"/power", messaging.Type(), bytes.NewReader(payload))
		if err != nil || resp.StatusCode != http.StatusOK {
			logger.Printf("Failed to kill: %v\n", err)
			time.Sleep(ACTION_DELAY_SEC * time.Second)
			continue
		}

		time.Sleep(ACTION_DELAY_SEC * time.Second)
	}

	g.POST("", postPower)
}
