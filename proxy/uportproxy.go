package proxy

import (
	"net"
	"shadowproxy/fillter"
	"shadowproxy/ids"
	"shadowproxy/logger"
	"sync"
	"time"
)

// UDP Port Proxy

type UDPConn struct {
	Addr     string
	Conn     net.Conn
	TTL      int64
	RecvTime time.Time
}

var UDPConns = map[string]*UDPConn{}
var UDPMutex = new(sync.Mutex)

func CleanTimeoutUDPConn() {

	for {
		UDPMutex.Lock()
		for k, v := range UDPConns {
			if time.Now().Sub(v.RecvTime).Milliseconds() > v.TTL {
				v.Conn.Close()
				delete(UDPConns, k)
			}
		}
		UDPMutex.Unlock()
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

}

func CleanAllUDPConn() {

	UDPMutex.Lock()
	defer UDPMutex.Unlock()

	for k, v := range UDPConns {
		v.Conn.Close()
		delete(UDPConns, k)
	}

}

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
		UDPMutex.Lock()
		defer UDPMutex.Unlock()

		udpConn, ok := UDPConns[addr.String()]

		if !ok {

			if fillter.Fillter(addr.String()) {

				logger.Warn("UDP", addr.String(), "Alice is filtrated")
			} else {

				ids.CheckAddr(addr.String())

				go UConnectionHandler(addr, listener, buffer, n1, backendAddr)
			}
			continue
		}

		ids.PackageLengthRecorder(addr.String(), n1)

		n2, err := UDPConns[addr.String()].Conn.Write(buffer[:n1])

		if err != nil {

			logger.Error("UDP", err)
			UDPConns[addr.String()].Conn.Close()
			delete(UDPConns, addr.String())
			continue
		}

		logger.Log("UDP", addr.String(), "->", udpConn.Conn.RemoteAddr().String(), n2, "Bytes")
		UDPConns[addr.String()].RecvTime = time.Now()
	}
}

func UConnectionHandler(addr *net.UDPAddr, listener *net.UDPConn, buffer []byte, n int, backendAddr string) {

	logger.Log("UDP", addr.String(), "Alice connected.")

	backend, err := net.Dial("udp", backendAddr)
	udpConn := new(UDPConn)

	if err != nil {
		logger.Error("UDP", err)
		return
	}

	logger.Log("UDP", backendAddr, "Bob connected.")
	udpConn.Addr = addr.String()
	udpConn.Conn = backend
	udpConn.TTL = 10000
	udpConn.RecvTime = time.Now()

	UDPConns[addr.String()] = udpConn
	LAddrToRAddr[backend.LocalAddr().String()] = addr.String()

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
		udpConn.RecvTime = time.Now()
	}
}

func init() {
	go CleanTimeoutUDPConn()
}
