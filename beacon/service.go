package beacon

import (
	"net"
	"time"

	"github.com/kaz/flos/messaging"
)

func getNIFs() ([]net.Interface, error) {
	nifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	result := []net.Interface{}
	for _, nif := range nifs {
		if nif.Flags&net.FlagUp == 0 {
			continue
		}
		if nif.Flags&net.FlagMulticast == 0 {
			continue
		}
		if nif.Flags&net.FlagLoopback != 0 {
			continue
		}

		result = append(result, nif)
	}

	return result, nil
}
func getAddrs4(nif net.Interface) ([]net.IP, error) {
	addrs, err := nif.Addrs()
	if err != nil {
		return nil, err
	}

	result := []net.IP{}
	for _, addr := range addrs {
		var ip net.IP = nil

		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP.To4()
		case *net.IPAddr:
			ip = v.IP.To4()
		}

		if ip != nil {
			result = append(result, ip)
		}
	}

	return result, nil
}

func startSender() {
	nifs, err := getNIFs()
	if err != nil {
		logger.Printf("failed to get interfaces: %v\n", err)
		return
	}

	for _, nif := range nifs {
		addrs, err := getAddrs4(nif)
		if err != nil {
			logger.Printf("failed to get addrs: %v\n", err)
			continue
		}

		for _, addr := range addrs {
			laddr, err := net.ResolveUDPAddr("udp", addr.String()+":11239")
			if err != nil {
				logger.Printf("failed to resolve addr: %v\n", err)
				continue
			}

			go func(nif net.Interface, laddr net.UDPAddr) {
				for {
					ch := make(chan error)
					go sendBeacon(ch, &laddr)
					logger.Printf("sending (dev=%s)\n", nif.Name)
					logger.Printf("send failed (dev=%s): %v\n", nif.Name, <-ch)
					close(ch)
				}
			}(nif, *laddr)
		}
	}
}
func sendBeacon(ch chan error, laddr *net.UDPAddr) {
	raddr, err := net.ResolveUDPAddr("udp", UDP_ADDR)
	if err != nil {
		ch <- err
		return
	}

	conn, err := net.DialUDP("udp", laddr, raddr)
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

		time.Sleep(BEACON_CYCLE_SEC * time.Second)
	}
}

func startReceiver() {
	nifs, err := getNIFs()
	if err != nil {
		logger.Printf("failed to get interfaces: %v\n", err)
		return
	}

	for _, nif := range nifs {
		addrs, err := getAddrs4(nif)
		if err != nil {
			logger.Printf("failed to get addrs: %v\n", err)
			continue
		}
		if len(addrs) == 0 {
			continue
		}

		go func(nif net.Interface) {
			for {
				ch := make(chan error)
				go recvBeacon(ch, &nif)
				logger.Printf("receiving (dev=%s)\n", nif.Name)
				logger.Printf("receive failed (dev=%s): %v\n", nif.Name, <-ch)
				close(ch)
			}
		}(nif)
	}
}
func recvBeacon(ch chan error, nif *net.Interface) {
	raddr, err := net.ResolveUDPAddr("udp", UDP_ADDR)
	if err != nil {
		ch <- err
		return
	}

	listener, err := net.ListenMulticastUDP("udp", nif, raddr)
	if err != nil {
		ch <- err
		return
	}
	defer listener.Close()

	// CAUTION: payload size must be less than 512
	buffer := make([]byte, 512)

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
		if _, ok := nodes[remote.IP.String()]; !ok {
			logger.Printf("Detected new node: %v\n", remote.IP)
		}
		nodes[remote.IP.String()] = received
		mu.Unlock()
	}
}
