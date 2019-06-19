package beacon

import (
	"net"
	"time"

	"github.com/kaz/flos/messaging"
)

func sendBeacon(ch chan error) {
	conn, err := net.Dial("udp", UDP_ADDR)
	if err != nil {
		ch <- err
		return
	}
	defer conn.Close()

	for {
		payload, err := messaging.Encode(time.Now())
		if err != nil {
			ch <- err
			return
		}
		if _, err := conn.Write(payload); err != nil {
			ch <- err
			return
		}
		logger.Println("Sent beacon")

		time.Sleep(BEACON_CYCLE_SEC * time.Second)
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
		n, remote, err := listener.ReadFromUDP(buffer)
		if err != nil {
			ch <- err
			return
		}

		var received time.Time
		if err := messaging.Decode(buffer[:n], &received); err != nil {
			logger.Printf("Beacon discarded: %v\n", err)
			continue
		}

		mu.Lock()
		nodes[remote.IP.String()] = received
		mu.Unlock()

		logger.Printf("Received beacon from %v\n", remote.IP)
	}
}
