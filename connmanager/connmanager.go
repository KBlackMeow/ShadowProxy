package connmanager

import (
	"net"
	"strings"
	"sync"
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

func CloseConnFromIP(addr string) {

	ip := strings.Split(addr, ":")[0]
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
