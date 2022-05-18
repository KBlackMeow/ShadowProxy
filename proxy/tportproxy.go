package proxy

import (
	"net"
	"shadowproxy/fillter"
	"shadowproxy/ids"
	"shadowproxy/logger"
	"time"
)

type IPConns struct {
	IP    string
	Conns []net.Conn
}

var IPToConns = map[string]*IPConns{}

func AddConnToIP(conn net.Conn, addr string) {
	ip := net.ParseIP(addr).String()

	ipConns, ok := IPToConns[ip]
	if ok {
		ipConns.Conns = append(ipConns.Conns, conn)
		return
	}
	IPToConns[addr] = &IPConns{IP: addr, Conns: []net.Conn{conn}}
}

func TimeoutCloseConn(addr string, dely uint64) {
	addr = net.ParseIP(addr).String()
	time.Sleep(time.Duration(dely) * time.Millisecond)

	ipConns, ok := IPToConns[addr]
	if ok {
		for _, v := range ipConns.Conns {
			v.Close()
		}
		delete(IPToConns, addr)
	}

}

// TCP Port Proxy
func RunTPortProxy(bindAddr string, backendAddr string) {
	listener, err := net.Listen("tcp4", bindAddr)

	if err != nil {
		logger.Error("TCP", err)
		WG.Done()
		return
	}

	defer listener.Close()
	logger.Log("TCP", bindAddr, "tcp-proxy started.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("TCP", err)
			WG.Done()
			return
		}
		ids.CheckIP(conn.RemoteAddr().String())
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
	defer backend.Close()

	SetRAddrToLAddr(backend.LocalAddr().String(), conn.RemoteAddr().String())
	AddConnToIP(backend, conn.RemoteAddr().String())
	logger.Log("TCP", backendAddr, "Bob connected.")

	closed := make(chan bool, 2)
	go TProxy(conn, backend, closed, true)
	go TProxy(backend, conn, closed, false)
	<-closed

	delete(LAddrToRAddr, backend.LocalAddr().String())
	delete(IPToConns, conn.RemoteAddr().String())

	logger.Log("TCP", conn.RemoteAddr().String(), "Alice connection is closed.")
}

func TProxy(from net.Conn, to net.Conn, closed chan bool, remhost bool) {
	buffer := make([]byte, 4096)

	for {
		n1, err := from.Read(buffer)
		if err != nil {
			closed <- true
			return
		}

		if remhost {
			ids.PackageLengthRecorder(from.RemoteAddr().String(), n1)
		}

		n2, err := to.Write(buffer[:n1])
		logger.Log("TCP", from.RemoteAddr().String(), "->", to.RemoteAddr().String(), n2, "Bytes")

		if err != nil {
			closed <- true
			return
		}
	}
}
