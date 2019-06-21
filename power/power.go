package power

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/kaz/flos/messaging"
	"github.com/labstack/echo/v4"
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

		payload, err := messaging.Encode("stop")
		if err != nil {
			logger.Printf("Failed to encode payload: %v\n", err)
			time.Sleep(ACTION_DELAY_SEC * time.Second)
			continue
		}

		logger.Println("Killing old process ...")
		resp, err := http.Post("http://"+LISTEN+"/power", messaging.Type(), bytes.NewReader(payload))
		if err != nil || resp.StatusCode != http.StatusOK {
			logger.Printf("Failed to kill: %v\n", err)
			time.Sleep(ACTION_DELAY_SEC * time.Second)
			continue
		}

		time.Sleep(ACTION_DELAY_SEC * time.Second)
	}
	go func() {
		for _, val := range os.Args {
			if val == "-d" || val == "--detach" {
				logger.Println("Detaching ...")
				restart()
			}
		}
	}()

	g.POST("", postPower)
}
