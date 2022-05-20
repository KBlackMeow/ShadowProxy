package proxy

import (
	"shadowproxy/config"
	"strings"
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

	for _, v := range config.ShadowProxyConfig.Rules {
		WG.Add(1)
		args := strings.Split(v, "://")
		if args[0] == "tcp" {
			go RunTPortProxy(args[1])
		} else if args[0] == "udp" {
			go RunUPortProxy(args[1])
		}
	}

	WG.Wait()

}
