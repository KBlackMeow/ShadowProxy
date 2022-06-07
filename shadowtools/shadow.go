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
var ShadowAddrs []string

func InitShadowService() {

	for _, serviceName := range config.ShadowProxyConfig.Shadows {
		serviceAddr, ok := transform.NameToAddr[serviceName]
		if ok {
			ShadowAddrs = append(ShadowAddrs, serviceAddr)
			continue

		}
		logger.Log(serviceName)
		addrs := strings.Split(serviceName, ":")
		addr := net.ParseIP(addrs[0])

		port, err := strconv.ParseInt(addrs[1], 10, 32)
		if addr != nil && err == nil && port < 65536 && port > 0 {
			ShadowAddrs = append(ShadowAddrs, serviceName)
			continue
		}
	}

}

func GetShadowAddr(remoteAddr string) string {
	port, err := strconv.ParseInt(strings.Split(remoteAddr, ":")[1], 10, 32)

	if err != nil {
		logger.Error(err)
		return ""
	}

	return ShadowAddrs[int(port)%len(ShadowAddrs)]
}
