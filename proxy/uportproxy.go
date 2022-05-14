package proxy

import (
	"net"
	"shadowproxy/fillter"
	"shadowproxy/ids"
	"shadowproxy/logger"
	"time"
)

// UDP Port Proxy

func RunUPortProxy(listenAddr, backendAddr string) {

	udpLAddr, _ := net.ResolveUDPAddr("udp", listenAddr)
	listener, err := net.ListenUDP("udp", udpLAddr)

	if err != nil {

		logger.Error("UDP", err)
		WG.Done()
		return
	}

	defer listener.Close()

	logger.Log("UDP", listenAddr, "udp-proxy started.")

	for {

		buffer := make([]byte, 4096)
		n1, addr, err := listener.ReadFromUDP(buffer)

		if err != nil {
			logger.Error("UDP", err)
			WG.Done()
			return
		}

		link, ok := links[addr.String()]

		if !ok {

			if fillter.Fillter(addr.String()) {

				logger.Warn("UDP", addr.String(), "Alice is filtrated")
			} else {

				ids.CheckAddr(addr.String())

				go UConnectionHandler(addr, listener, buffer, n1, backendAddr, links)
			}
			continue
		}

		ids.PackageLengthRecorder(addr.String(), n1)

		n2, err := links[addr.String()].backend.Write(buffer[:n1])

		if err != nil {

			logger.Error("UDP", err)
			links[addr.String()].backend.Close()
			delete(links, addr.String())
			continue
		}

		logger.Log("UDP", addr.String(), "->", link.backend.RemoteAddr().String(), n2, "Bytes")
		links[addr.String()].recvtime = time.Now()
	}
}

func UConnectionHandler(addr *net.UDPAddr, listener *net.UDPConn, buffer []byte, n int, backendAddr string, links map[string]*UDPLink) {

	logger.Log("UDP", addr.String(), "Alice connected.")

	backend, err := net.Dial("udp", backendAddr)
	link := new(UDPLink)

	if err != nil {
		logger.Error("UDP", err)
		return
	}

	logger.Log("UDP", backendAddr, "Bob connected.")
	link.addr = addr.String()
	link.backend = backend
	link.ttl = 10000
	link.recvtime = time.Now()

	links[addr.String()] = link

	n2, err := backend.Write(buffer[:n])
	if err != nil {
		logger.Error("UDP", err)
		return
	}

	logger.Log("UDP", addr.String(), "->", backendAddr, n2, "Bytes")

	defer backend.Close()

	for {

		buffer := make([]byte, 4096)
		n1, err := backend.Read(buffer)

		if err != nil {
			logger.Error("UDP", err)
			return
		}

		n2, err := listener.WriteToUDP(buffer[:n1], addr)

		if err != nil {
			logger.Error("UDP", err)
			return
		}

		logger.Log("UDP", backendAddr, "->", addr.String(), n2, "Bytes")
		link.recvtime = time.Now()
	}
}
