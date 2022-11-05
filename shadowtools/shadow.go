package shadowtools

import (
	"net"
	"shadowproxy/config"
	"shadowproxy/logger"
	"strconv"
	"strings"
)

// var ShadowAddr string
var ShadowAddr string

func InitShadowService() {

	serviceAddr := config.ShadowProxyConfig.Shadow

	logger.Log("Shadow Addr:", serviceAddr)
	addrs := strings.Split(serviceAddr, ":")
	addr := net.ParseIP(addrs[0])

	port, err := strconv.ParseInt(addrs[1], 10, 32)
	if addr != nil && err == nil && port < 65536 && port > 0 {
		ShadowAddr = serviceAddr
		return
	}

}

func GetShadowAddr() string {

	return ShadowAddr
}
