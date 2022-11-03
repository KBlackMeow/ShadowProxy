package proxy

import (
	"net"
	"shadowproxy/filter"
	"shadowproxy/ids"
	"shadowproxy/logger"
	"shadowproxy/transform"
	"strings"
	"sync"
	"time"
)

// UDP Port Proxy

type UDPConn struct {
	addr         *net.UDPAddr
	backendConn  net.Conn
	listenerConn *net.UDPConn
	TTL          int64
	RecvTime     time.Time
}

var UDPConns = map[string]*UDPConn{}
var UDPMutex = new(sync.Mutex)

func SetUDPConn(addr string, conn *UDPConn) {

	UDPMutex.Lock()
	defer UDPMutex.Unlock()
	UDPConns[addr] = conn

}

func GetUDPConn(addr string) (*UDPConn, bool) {

	UDPMutex.Lock()
	defer UDPMutex.Unlock()
	udpConn, ok := UDPConns[addr]
	return udpConn, ok

}

func (conn UDPConn) Write(buff []byte) (int, error) {

	return conn.listenerConn.WriteToUDP(buff, conn.addr)

}

func (conn UDPConn) Close() {

	conn.backendConn.Close()

}

func UDPConnClear() {

	for {
		for k, v := range UDPConns {
			if time.Now().Sub(v.RecvTime).Milliseconds() > v.TTL {
				v.Close()
				delete(UDPConns, k)
			}
		}

		time.Sleep(time.Duration(500) * time.Millisecond)
	}

}

func AllUDPConnClear() {

	for k, v := range UDPConns {
		v.Close()
		delete(UDPConns, k)
	}

}

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
	proxy.forword()

}
func (proxy UDProxy) forword() {

	for {
		buffer := make([]byte, 4096)
		n1, addr, err := proxy.listener.ReadFromUDP(buffer)

		if err != nil {
			logger.Error("UDP", err)
			WG.Done()
			return
		}

		udpConn, ok := GetUDPConn(addr.String())

		if !ok {
			if filter.Filter(addr.String()) {
				logger.Warn("UDP", addr.String(), "Alice is filtrated")
				continue
			} else {
				ids.CheckIP(addr.String())
				udpConn = link(proxy.listener, addr, proxy.backendAddr)
			}

		}

		ids.PackageLengthRecorder(addr.String(), n1)
		n2, err := udpConn.backendConn.Write(buffer[:n1])

		if err != nil {
			logger.Error("UDP", err)
			UDPConns[addr.String()].backendConn.Close()
			delete(UDPConns, addr.String())
			continue
		}

		logger.Log("UDP", addr.String(), "->", udpConn.backendConn.RemoteAddr().String(), n2, "Bytes")
		UDPConns[addr.String()].RecvTime = time.Now()
	}
}

func link(listener *net.UDPConn, addr *net.UDPAddr, backendAddr string) *UDPConn {

	logger.Log("UDP", addr.String(), "Alice connected.")
	backend, err := net.Dial("udp", backendAddr)
	if err != nil {
		logger.Error("UDP", err)
		return nil
	}
	// defer backend.Close()

	logger.Log("UDP", backendAddr, "Bob connected.")

	conn := new(UDPConn)
	conn.addr = addr
	conn.backendConn = backend
	conn.TTL = 10000
	conn.RecvTime = time.Now()
	conn.listenerConn = listener

	SetUDPConn(addr.String(), conn)
	transform.SetRemoteAddrToLocalAddr(backend.LocalAddr().String(), addr.String())

	go conn.back()
	return conn

}

func (conn *UDPConn) back() {
	from := conn.backendConn
	to := conn
	for {
		buffer := make([]byte, 4096)
		n1, err := from.Read(buffer)

		if err != nil {
			logger.Error("UDP", err)
			conn.Close()
			return
		}

		n2, err := to.Write(buffer[:n1])
		if err != nil {
			logger.Error("UDP", err)
			conn.Close()
			return
		}

		logger.Log("UDP", from.RemoteAddr().String(), "->", to.addr.String(), n2, "Bytes")
		to.RecvTime = time.Now()
	}
}

func RunUPortProxy(rule string) {

	args := strings.Split(rule, "->")
	if len(args) == 2 {
		proxy := UDProxy{bindAddr: args[0], backendAddr: args[1]}
		go proxy.Run()
	}

}

func init() {

	go UDPConnClear()

}
