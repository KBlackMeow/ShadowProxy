package proxy

import (
	"shadowproxy/config"
)

var LAddrToRAddr = map[string]string{}
var NameToAddr = map[string]string{}
var ShadowAddr = ""

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
		return
	}

	WG.Wait()
}
