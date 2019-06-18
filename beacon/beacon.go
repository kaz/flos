package beacon

import (
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/kaz/flos/messaging"
)

const (
	BEACON_CYCLE = 5
	UDP_ADDR     = "239.239.239.239:239"
)

var (
	logger = log.New(os.Stdout, "[beacon] ", log.Ltime)

	mu    = sync.RWMutex{}
	nodes = map[string]time.Time{}
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
		nodes[string(remote.IP)] = received
		mu.Unlock()

		logger.Printf("Received beacon from %v\n", remote.IP)
	}
}

func GetNodes() map[string]time.Time {
	mu.RLock()
	defer mu.RUnlock()

	result := make(map[string]time.Time)
	for k, v := range nodes {
		result[k] = v
	}

	return result
}
func DeleteNode(ip net.IP) {
	mu.Lock()
	defer mu.Unlock()

	delete(nodes, string(ip))
}
