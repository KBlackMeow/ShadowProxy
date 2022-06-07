package proxy

import (
	"shadowproxy/config"
	"strings"
	"sync"
)

var WG sync.WaitGroup

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
