package shadowtools

import (
	"net"
	"shadowproxy/config"
	"shadowproxy/logger"
	"shadowproxy/transform"
	"strconv"
	"strings"
)

// var ShadowAddr string
var ShadowAddr string

func InitShadowService() {

	serviceName := config.ShadowProxyConfig.Shadow
	serviceAddr, ok := transform.NameToAddr[serviceName]
	if ok {
		ShadowAddr = serviceAddr
		return

	}
	logger.Log(serviceName)
	addrs := strings.Split(serviceName, ":")
	addr := net.ParseIP(addrs[0])

	port, err := strconv.ParseInt(addrs[1], 10, 32)
	if addr != nil && err == nil && port < 65536 && port > 0 {
		ShadowAddr = serviceName
		return
	}

}

func GetShadowAddr() string {

	return ShadowAddr
}
