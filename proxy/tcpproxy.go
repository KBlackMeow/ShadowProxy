package proxy

import (
	"net"
	"shadowproxy/fillter"
	"shadowproxy/ids"
	"shadowproxy/logger"
	"shadowproxy/shadowtools"
	"shadowproxy/transform"
	"strings"
	"sync"
	"time"
)

type IPConns struct {
	IP    string
	Conns []net.Conn
}

var IPToConns = map[string]*IPConns{}
var TCPMutex = new(sync.Mutex)

func AddConnToIP(conn net.Conn, addr string) {

	TCPMutex.Lock()
	defer TCPMutex.Unlock()
	ip := strings.Split(addr, ":")[0]

	ipConns, ok := IPToConns[ip]
	if ok {
		ipConns.Conns = append(ipConns.Conns, conn)
		return
	}
	IPToConns[ip] = &IPConns{IP: addr, Conns: []net.Conn{conn}}

}

func ClearConnFromIP(addr string, dely uint64) {

	ip := strings.Split(addr, ":")[0]
	time.Sleep(time.Duration(dely) * time.Millisecond)
	TCPMutex.Lock()
	defer TCPMutex.Unlock()
	ipConns, ok := IPToConns[ip]
	if ok {
		for _, v := range ipConns.Conns {
			v.Close()
		}
		delete(IPToConns, addr)
	}

}

type TCPProxy struct {
	bindAddr    string
	backendAddr string
}

// TCP Port Proxy
func (proxy TCPProxy) Run() {

	listener, err := net.Listen("tcp4", proxy.bindAddr)

	if err != nil {
		logger.Error("TCP", err)
		WG.Done()
		return
	}

	defer listener.Close()
	logger.Log("TCP", proxy.bindAddr, "tcp-proxy started.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("TCP", err)
			WG.Done()
			return
		}
		ids.CheckIP(conn.RemoteAddr().String())

		shadowAddr := shadowtools.GetShadowAddr()

		if fillter.Fillter(conn.RemoteAddr().String()) {
			logger.Warn("TCP", conn.RemoteAddr().String(), "Alice is filtrated", "Shadow is", shadowAddr)
			go handler(conn, shadowAddr)

		} else {
			go handler(conn, proxy.backendAddr)
		}
	}

}

func handler(conn net.Conn, backendAddr string) {

	if backendAddr == "" {
		conn.Close()
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

	transform.SetRAddrToLAddr(backend.LocalAddr().String(), conn.RemoteAddr().String())
	AddConnToIP(backend, conn.RemoteAddr().String())
	logger.Log("TCP", backendAddr, "Bob connected.")

	closed := make(chan bool, 2)
	go proxy(conn, backend, closed, true)
	go proxy(backend, conn, closed, false)
	<-closed

	transform.DeleteAddr(backend.LocalAddr().String())

	delete(IPToConns, conn.RemoteAddr().String())

	logger.Log("TCP", conn.RemoteAddr().String(), "Alice connection is closed.")

}

func proxy(from net.Conn, to net.Conn, closed chan bool, RTL bool) {

	buffer := make([]byte, 4096)

	for {
		n1, err := from.Read(buffer)
		if err != nil {
			closed <- true
			return
		}

		if RTL {
			ids.PackageLengthRecorder(from.RemoteAddr().String(), n1)
		} else {
		}

		n2, err := to.Write(buffer[:n1])
		logger.Log("TCP", from.RemoteAddr().String(), "->", to.RemoteAddr().String(), n2, "Bytes")

		if err != nil {
			closed <- true
			return
		}
	}

}

func RunTPortProxy(rule string) {

	args := strings.Split(rule, "->")
	if len(args) == 2 {
		proxy := TCPProxy{bindAddr: args[0], backendAddr: args[1]}
		go proxy.Run()
	}

}
