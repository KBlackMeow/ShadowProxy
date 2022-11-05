package proxy

import (
	"net"
	"shadowproxy/connmanager"
	"shadowproxy/filter"
	"shadowproxy/ids"
	"shadowproxy/logger"
	"shadowproxy/transform"
	"strings"
	"time"
)

// UDP Port Proxy

type UDProxy struct {
	bindAddr    string
	backendAddr string
	listener    *net.UDPConn
}

func (proxy UDProxy) Run() {

	udpLAddr, _ := net.ResolveUDPAddr("udp4", proxy.bindAddr)
	listener, err := net.ListenUDP("udp4", udpLAddr)

	if err != nil {
		logger.Error("UDP", err)
		WG.Done()
		return
	}

	defer listener.Close()
	logger.Log("UDP", proxy.bindAddr, "udp-proxy started.")
	proxy.listener = listener
	proxy.Forword()

}
func (proxy UDProxy) Forword() {

	for {
		buffer := make([]byte, 4096)
		n1, addr, err := proxy.listener.ReadFromUDP(buffer)

		if err != nil {
			logger.Error("UDP", err)
			WG.Done()
			return
		}

		if filter.Filter(addr.String()) {
			logger.Warn("UDP", addr.String(), "Alice is filtrated")
			continue
		} else {
			udpConn, ok := connmanager.GetUDPConn(addr.String())
			if !ok {
				ids.CheckIP(addr.String())
				udpConn = link(proxy.listener, addr, proxy.backendAddr)
			}

			ids.PackageLengthRecorder(addr.String(), n1)
			n2, err := udpConn.BackendConn.Write(buffer[:n1])

			if err != nil {
				logger.Error("UDP", err)
				connmanager.CloseUDPConnFromAddr(addr.String())
				continue
			}

			logger.Log("UDP", addr.String(), "->", udpConn.BackendConn.RemoteAddr().String(), n2, "Bytes")
			connmanager.UDPConns[addr.String()].RecvTime = time.Now()

		}

	}
}

func link(listener *net.UDPConn, addr *net.UDPAddr, backendAddr string) *connmanager.UDPConn {

	logger.Log("UDP", addr.String(), "Alice connected.")
	backend, err := net.Dial("udp", backendAddr)
	if err != nil {
		logger.Error("UDP", err)
		return nil
	}
	// defer backend.Close()

	logger.Log("UDP", backendAddr, "Bob connected.")

	conn := new(connmanager.UDPConn)
	conn.Addr = addr
	conn.BackendConn = backend
	conn.TTL = 10000
	conn.RecvTime = time.Now()
	conn.ListenerConn = listener

	connmanager.SetUDPConn(addr.String(), conn)
	transform.SetRemoteAddrToLocalAddr(backend.LocalAddr().String(), addr.String())

	go conn.Backword()
	return conn

}

func RunUPortProxy(rule string) {

	args := strings.Split(rule, "->")
	if len(args) == 2 {
		proxy := UDProxy{bindAddr: args[0], backendAddr: args[1]}
		go proxy.Run()
	}

}
