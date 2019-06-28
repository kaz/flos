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
	DELAY_SEC = 2
	LISTEN    = ":10239"
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
		} else {
			logger.Println("Killing old process ...")
			resp, err := http.Post("http://"+LISTEN+"/power", messaging.Type(), bytes.NewReader(payload))
			if err != nil {
				logger.Printf("Failed to kill: %v\n", err)
			} else if resp.StatusCode != http.StatusOK {
				logger.Printf("Failed to kill: (http_status=%d)\n", resp.StatusCode)
			}
		}
		time.Sleep(DELAY_SEC * time.Second)
	}
	go func() {
		time.Sleep(DELAY_SEC * time.Second)

		for _, val := range os.Args {
			if val == "-d" || val == "--detach" {
				logger.Println("Detaching ...")
				restart()
			}
		}
	}()

	g.POST("", postPower)
}
