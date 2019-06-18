package beacon

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/kaz/flos/messaging"
	"github.com/labstack/echo/v4"
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

func StartService(g *echo.Group) {
	g.GET("/nodes", getNodes)
	g.DELETE("/node", deleteNode)

	go send()
	go recv()
}

func send() {
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

func recv() {
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
		nodes[remote.IP.String()] = received
		mu.Unlock()

		logger.Printf("Received beacon from %v\n", remote.IP)
	}
}

func getNodes(c echo.Context) error {
	mu.RLock()
	defer mu.RUnlock()

	resp := make(map[string]time.Time)
	for k, v := range nodes {
		resp[k] = v
	}

	c.Set("response", resp)
	return nil
}
func deleteNode(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()

	req, ok := c.Get("request").(string)
	if !ok {
		return fmt.Errorf("unpected request format")
	}

	delete(nodes, req)
	return nil
}
