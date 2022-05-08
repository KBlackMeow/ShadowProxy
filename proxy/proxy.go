package proxy

import (
	"flag"
	"net"
	"sync"
	"time"
)

type UDPlink struct {
	addr     string
	backend  net.Conn
	ttl      uint64
	recvtime time.Time
}

var LAddrToRAddr = map[string]string{}
var NameToAddr = map[string]string{}
var ShadowAddr = ""
var links = map[string]*UDPlink{}
var AddrToConn = map[string]net.Conn{}
var WG sync.WaitGroup
var ProxyProtocol = "tcp"
var ProxyBindAddr = "0.0.0.0:30000"
var ProxyBackendAddr = "127.0.0.1:30000"

func CleanTimeoutConn() {

	for {
		for k, v := range links {
			if uint64(time.Now().Sub(v.recvtime).Nanoseconds()/1e6) > v.ttl {
				v.backend.Close()
				delete(links, k)
			}
		}
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

}

func CleanAllConn() {

	for k, v := range links {
		v.backend.Close()
		delete(links, k)
	}

}

func TimeoutCloseConn(addr string, dely uint64) {

	time.Sleep(time.Duration(dely) * time.Millisecond)
	conn, ok := AddrToConn[addr]
	if ok {
		conn.Close()
		// delete(AddrToConn, addr)
	}

}

func init() {
	go CleanTimeoutConn()
}

func RunProxy() {
	if ProxyProtocol == "tcp" {
		WG.Add(1)
		go RunTPortProxy(ProxyBindAddr, ProxyBackendAddr)

	} else if ProxyProtocol == "udp" {
		WG.Add(1)
		go RunUPortProxy(ProxyBindAddr, ProxyBackendAddr)

	} else if ProxyProtocol == "tcp/udp" {
		WG.Add(2)
		go RunTPortProxy(ProxyBindAddr, ProxyBackendAddr)
		go RunUPortProxy(ProxyBindAddr, ProxyBackendAddr)

	} else {
		flag.Usage()
		return
	}

	WG.Wait()
}
