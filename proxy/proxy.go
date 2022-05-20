package proxy

import (
	"shadowproxy/config"
	"sync"
)

var WG sync.WaitGroup
var Mutex = new(sync.Mutex)

var LAddrToRAddr = map[string]string{}

func GetRAddrFromLAddr(laddr string) (string, bool) {

	Mutex.Lock()
	defer Mutex.Unlock()
	raddr, ok := LAddrToRAddr[laddr]
	return raddr, ok

}

func SetRAddrToLAddr(laddr string, raddr string) {

	Mutex.Lock()
	defer Mutex.Unlock()
	LAddrToRAddr[laddr] = raddr

}

func RunProxy() {

	if config.ShadowProxyConfig.Protocol == "tcp" {
		WG.Add(1)
		go RunTPortProxy(config.ShadowProxyConfig.Rules)

	} else if config.ShadowProxyConfig.Protocol == "udp" {
		WG.Add(1)
		go RunUPortProxy(config.ShadowProxyConfig.Rules)

	} else if config.ShadowProxyConfig.Protocol == "tcp/udp" {
		WG.Add(2)
		go RunTPortProxy(config.ShadowProxyConfig.Rules)
		go RunUPortProxy(config.ShadowProxyConfig.Rules)

	} else {
		return
	}

	WG.Wait()

}
