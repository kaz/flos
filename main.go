package main

import (
	"time"

	"github.com/kaz/flos/beacon"
)

func main() {
	go beacon.RecvBeacon()
	go beacon.SendBeacon()

	time.Sleep(12 * time.Hour)
}
