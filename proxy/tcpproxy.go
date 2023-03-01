package proxy

import (
	"net"
	"shadowproxy/connmanager"
	"shadowproxy/filter"
	"shadowproxy/ids"
	"shadowproxy/logger"
	"shadowproxy/shadowtools"
	"shadowproxy/transform"
	"strings"
)

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
	logger.Log("TCP", proxy.bindAddr, "->", proxy.backendAddr, "tcp-proxy started.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("TCP", err)
			WG.Done()
			return
		}
		go ids.CheckIP(conn.RemoteAddr().String())

		shadowAddr := shadowtools.GetShadowAddr()

		if filter.Filter(conn.RemoteAddr().String()) {
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

	transform.SetRemoteAddrToLocalAddr(backend.LocalAddr().String(), conn.RemoteAddr().String())
	connmanager.AddConnToIP(backend, conn.RemoteAddr().String())
	logger.Log("TCP", backendAddr, "Bob connected.")

	closed := make(chan bool, 2)
	go proxy(conn, backend, closed, true)
	go proxy(backend, conn, closed, false)
	<-closed

	transform.DeleteAddr(backend.LocalAddr().String())

	connmanager.CloseConnFromIP(conn.RemoteAddr().String())
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
			go ids.PackageLengthRecorder(from.RemoteAddr().String(), n1)
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
