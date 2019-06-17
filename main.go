package main

import (
	"time"

	"github.com/kaz/flosd/beacon"
)

func main() {
	go beacon.RecvBeacon()
	go beacon.SendBeacon()

	time.Sleep(12 * time.Hour)
}
