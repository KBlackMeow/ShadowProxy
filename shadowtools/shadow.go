package shadowtools

import (
	"net"
	"shadowproxy/config"
	"shadowproxy/proxy"
	"strconv"
	"strings"
)

var ShadowAddr string

func InitShadowService() {
	serviceName := config.ShadowProxyConfig.Shadow

	serviceAddr, ok := proxy.NameToAddr[serviceName]
	if ok {
		proxy.ShadowAddr = serviceAddr
		return
	}

	addrs := strings.Split(serviceName, ":")
	addr := net.ParseIP(addrs[0])
	port, err := strconv.ParseInt(addrs[1], 10, 32)

	// logger.Log(addr, port)
	if addr != nil && err == nil && port < 65536 && port > 0 {
		proxy.ShadowAddr = serviceName
		return
	}
	proxy.ShadowAddr = ""
}
