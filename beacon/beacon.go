package beacon

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/kaz/flosd/messaging"
)

var logger = log.New(os.Stdout, "[beacon] ", log.Ltime)

const (
	BEACON_CYCLE = 5
	UDP_ADDR     = "239.239.239.239:239"
	PAYLOAD      = "*** FLOS ***"
)

func SendBeacon() {
	conn, err := net.Dial("udp", UDP_ADDR)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		payload, err := messaging.Encode(PAYLOAD)
		if err != nil {
			panic(err)
		}
		if _, err := conn.Write(payload); err != nil {
			panic(err)
		}
		logger.Println("Sent beacon")

		time.Sleep(BEACON_CYCLE * time.Second)
	}
}

func RecvBeacon() {
	address, err := net.ResolveUDPAddr("udp", UDP_ADDR)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenMulticastUDP("udp", nil, address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	buffer := make([]byte, 256*1024)
	for {
		n, remoteAddress, err := listener.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}

		data, err := messaging.Decode(buffer[:n])
		if err != nil {
			logger.Printf("Ignored: %v\n", err)
			continue
		}

		if data != PAYLOAD {
			logger.Printf("Ignored: invalid payload")
			continue
		}

		logger.Printf("Received from %v\n", remoteAddress.IP)
	}
}
