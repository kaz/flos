package beacon

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/kaz/flos/messaging"
)

var logger = log.New(os.Stdout, "[beacon] ", log.Ltime)

const (
	BEACON_CYCLE = 5
	UDP_ADDR     = "239.239.239.239:239"
	PAYLOAD      = "*** FLOS ***"
)

func SendBeacon() {
	for {
		ch := make(chan error)
		go sendBeacon(ch)
		logger.Printf("Sending beacon failed: %v\n", <-ch)
		close(ch)
	}
}
func sendBeacon(ch chan error) {
	conn, err := net.Dial("udp", UDP_ADDR)
	if err != nil {
		ch <- err
		return
	}
	defer conn.Close()

	for {
		payload, err := messaging.Encode(PAYLOAD)
		if err != nil {
			ch <- err
			return
		}
		if _, err := conn.Write(payload); err != nil {
			ch <- err
			return
		}
		logger.Println("Sent beacon")

		time.Sleep(BEACON_CYCLE * time.Second)
	}
}

func RecvBeacon() {
	for {
		ch := make(chan error)
		go recvBeacon(ch)
		logger.Printf("Receiving beacon failed: %v\n", <-ch)
		close(ch)
	}
}
func recvBeacon(ch chan error) {
	address, err := net.ResolveUDPAddr("udp", UDP_ADDR)
	if err != nil {
		ch <- err
		return
	}

	listener, err := net.ListenMulticastUDP("udp", nil, address)
	if err != nil {
		ch <- err
		return
	}
	defer listener.Close()

	buffer := make([]byte, 256*1024)
	for {
		n, remoteAddress, err := listener.ReadFromUDP(buffer)
		if err != nil {
			ch <- err
			return
		}

		var payload string
		if err := messaging.Decode(buffer[:n], &payload); err != nil {
			logger.Printf("Ignored: %v\n", err)
			continue
		}

		if payload != PAYLOAD {
			logger.Printf("Ignored: invalid payload")
			continue
		}

		logger.Printf("Received beacon from %v\n", remoteAddress.IP)
	}
}
