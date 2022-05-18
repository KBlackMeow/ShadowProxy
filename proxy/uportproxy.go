package proxy

import (
	"net"
	"shadowproxy/fillter"
	"shadowproxy/ids"
	"shadowproxy/logger"
	"time"
)

// UDP Port Proxy

func RunUPortProxy(bindAddr, backendAddr string) {

	udpLAddr, _ := net.ResolveUDPAddr("udp", bindAddr)
	listener, err := net.ListenUDP("udp", udpLAddr)

	if err != nil {

		logger.Error("UDP", err)
		WG.Done()
		return
	}

	defer listener.Close()

	logger.Log("UDP", bindAddr, "udp-proxy started.")

	for {

		buffer := make([]byte, 4096)
		n1, addr, err := listener.ReadFromUDP(buffer)

		if err != nil {
			logger.Error("UDP", err)
			WG.Done()
			return
		}

		conn, ok := UDPConns[addr.String()]

		if !ok {

			if fillter.Fillter(addr.String()) {

				logger.Warn("UDP", addr.String(), "Alice is filtrated")
			} else {

				ids.CheckAddr(addr.String())

				go UConnectionHandler(addr, listener, buffer, n1, backendAddr, UDPConns)
			}
			continue
		}

		ids.PackageLengthRecorder(addr.String(), n1)

		n2, err := UDPConns[addr.String()].backend.Write(buffer[:n1])

		if err != nil {

			logger.Error("UDP", err)
			UDPConns[addr.String()].backend.Close()
			delete(UDPConns, addr.String())
			continue
		}

		logger.Log("UDP", addr.String(), "->", conn.backend.RemoteAddr().String(), n2, "Bytes")
		UDPConns[addr.String()].recvtime = time.Now()
	}
}

func UConnectionHandler(addr *net.UDPAddr, listener *net.UDPConn, buffer []byte, n int, backendAddr string, conns map[string]*UDPConn) {

	logger.Log("UDP", addr.String(), "Alice connected.")

	backend, err := net.Dial("udp", backendAddr)
	conn := new(UDPConn)

	if err != nil {
		logger.Error("UDP", err)
		return
	}

	logger.Log("UDP", backendAddr, "Bob connected.")
	conn.addr = addr.String()
	conn.backend = backend
	conn.ttl = 10000
	conn.recvtime = time.Now()

	conns[addr.String()] = conn

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
		conn.recvtime = time.Now()
	}
}
