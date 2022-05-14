package proxy

import (
	"flag"
	"net"
	"sync"
	"time"
)

type UDPLink struct {
	addr     string
	backend  net.Conn
	ttl      uint64
	recvtime time.Time
}

type IPLinks struct {
	IP    string
	links []net.Conn
}

var LAddrToRAddr = map[string]string{}
var NameToAddr = map[string]string{}
var ShadowAddr = ""
var links = map[string]*UDPLink{}
var IPToLinks = map[string]*IPLinks{}
var WG sync.WaitGroup
var ProxyProtocol = "tcp"
var ProxyBindAddr = "0.0.0.0:30000"
var ProxyBackendAddr = "127.0.0.1:30000"

var Mutex = new(sync.Mutex)

func AddLinkToIP(conn net.Conn, addr string) {

	addr = net.ParseIP(addr).String()

	Mutex.Lock()
	defer Mutex.Unlock()

	links, ok := IPToLinks[addr]
	if ok {
		links.links = append(links.links, conn)
		return
	}
	IPToLinks[addr] = &IPLinks{IP: addr, links: []net.Conn{conn}}

}

func CleanTimeoutConn() {

	for {
		Mutex.Lock()
		for k, v := range links {
			if uint64(time.Now().Sub(v.recvtime).Nanoseconds()/1e6) > v.ttl {
				v.backend.Close()
				delete(links, k)
			}
		}
		Mutex.Unlock()
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

}

func CleanAllConn() {

	Mutex.Lock()
	defer Mutex.Unlock()

	for k, v := range links {
		v.backend.Close()
		delete(links, k)
	}

}

func TimeoutCloseConn(addr string, dely uint64) {

	addr = net.ParseIP(addr).String()

	time.Sleep(time.Duration(dely) * time.Millisecond)

	Mutex.Lock()
	defer Mutex.Unlock()
	ipLinks, ok := IPToLinks[addr]
	if ok {
		for _, v := range ipLinks.links {
			v.Close()
		}
		delete(IPToLinks, addr)
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
