package proxy

import (
	"net"
	"shadowproxy/fillter"
	"shadowproxy/ids"
	"shadowproxy/logger"
)

// TCP Port Proxy
func RunTPortProxy(listenAddr string, backendAddr string) {

	listener, err := net.Listen("tcp4", listenAddr)

	if err != nil {
		logger.Error("TCP", err)
		WG.Done()
		return
	}

	defer listener.Close()

	logger.Log("TCP", listenAddr, "tcp-proxy started.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("TCP", err)
			WG.Done()
			return
		}
		ids.CheckAddr(conn.RemoteAddr().String())
		if fillter.Fillter(conn.RemoteAddr().String()) {

			logger.Warn("TCP", conn.RemoteAddr().String(), "Alice is filtrated", "Shadow is", ShadowAddr)
			go TConnectionHandler(conn, ShadowAddr)

		} else {

			go TConnectionHandler(conn, backendAddr)
		}
	}
}

func TConnectionHandler(conn net.Conn, backendAddr string) {

	if backendAddr == "" {

		conn.Close()
		logger.Log("TCP", conn.RemoteAddr().String(), "Alice is closed due to security strategy")
		return
	}

	logger.Log("TCP", conn.RemoteAddr().String(), "Alice connected.")
	backend, err := net.Dial("tcp", backendAddr)

	defer conn.Close()

	if err != nil {
		logger.Error("TCP", err)
		return
	}

	LAddrToRAddr[backend.LocalAddr().String()] = conn.RemoteAddr().String()

	AddLinkToIP(backend, conn.RemoteAddr().String())
	// logger.Log(LAddrToRAddr, AddrToConn)

	defer backend.Close()

	logger.Log("TCP", backendAddr, "Bob connected.")

	closed := make(chan bool, 2)

	go TProxy(conn, backend, closed)
	go TProxy(backend, conn, closed)
	<-closed

	delete(LAddrToRAddr, backend.LocalAddr().String())
	delete(IPToLinks, conn.RemoteAddr().String())
	// logger.Log(LAddrToRAddr, AddrToConn)
	logger.Log("TCP", conn.RemoteAddr().String(), "Alice connection is closed.")
}

func TProxy(from net.Conn, to net.Conn, closed chan bool) {

	buffer := make([]byte, 4096)
	for {

		n1, err := from.Read(buffer)
		if err != nil {

			closed <- true
			return
		}

		ids.PackageLengthRecorder(from.RemoteAddr().String(), n1)

		n2, err := to.Write(buffer[:n1])
		logger.Log("TCP", from.RemoteAddr().String(), "->", to.RemoteAddr().String(), n2, "Bytes")

		if err != nil {

			closed <- true
			return
		}
	}
}
