package proxy

import (
	"flag"
	"net"
	"shadowproxy/config"
	"sync"
	"time"
)

type UDPConn struct {
	addr     string
	backend  net.Conn
	ttl      uint64
	recvtime time.Time
}

type IPConns struct {
	IP    string
	Conns []net.Conn
}

var LAddrToRAddr = map[string]string{}
var NameToAddr = map[string]string{}
var ShadowAddr = ""
var UDPConns = map[string]*UDPConn{}
var IPToConns = map[string]*IPConns{}
var WG sync.WaitGroup

var Mutex = new(sync.Mutex)

func AddConnToIP(conn net.Conn, addr string) {

	addr = net.ParseIP(addr).String()

	Mutex.Lock()
	defer Mutex.Unlock()

	Conns, ok := IPToConns[addr]
	if ok {
		Conns.Conns = append(Conns.Conns, conn)
		return
	}
	IPToConns[addr] = &IPConns{IP: addr, Conns: []net.Conn{conn}}

}

func CleanTimeoutUDPConn() {

	for {
		Mutex.Lock()
		for k, v := range UDPConns {
			if uint64(time.Now().Sub(v.recvtime).Nanoseconds()/1e6) > v.ttl {
				v.backend.Close()
				delete(UDPConns, k)
			}
		}
		Mutex.Unlock()
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

}

func CleanAllUDPConn() {

	Mutex.Lock()
	defer Mutex.Unlock()

	for k, v := range UDPConns {
		v.backend.Close()
		delete(UDPConns, k)
	}

}

func TimeoutCloseConn(addr string, dely uint64) {

	addr = net.ParseIP(addr).String()

	time.Sleep(time.Duration(dely) * time.Millisecond)

	Mutex.Lock()
	defer Mutex.Unlock()
	ipConns, ok := IPToConns[addr]
	if ok {
		for _, v := range ipConns.Conns {
			v.Close()
		}
		delete(IPToConns, addr)
	}

}

func init() {
	go CleanTimeoutUDPConn()
}

func RunProxy() {
	if config.ShadowProxyConfig.Protocol == "tcp" {
		WG.Add(1)
		go RunTPortProxy(config.ShadowProxyConfig.BindAddr, config.ShadowProxyConfig.BackendAddr)

	} else if config.ShadowProxyConfig.Protocol == "udp" {
		WG.Add(1)
		go RunUPortProxy(config.ShadowProxyConfig.BindAddr, config.ShadowProxyConfig.BackendAddr)

	} else if config.ShadowProxyConfig.Protocol == "tcp/udp" {
		WG.Add(2)
		go RunTPortProxy(config.ShadowProxyConfig.BindAddr, config.ShadowProxyConfig.BackendAddr)
		go RunUPortProxy(config.ShadowProxyConfig.BindAddr, config.ShadowProxyConfig.BackendAddr)

	} else {
		flag.Usage()
		return
	}

	WG.Wait()
}
