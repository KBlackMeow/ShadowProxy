package shadowtools

import (
	"shadowproxy/proxy"
)

var ShadowAddr string

func SetShadowService(serviceName string) {

	serviceAddr, ok := proxy.NameToAddr[serviceName]
	if !ok {
		proxy.ShadowAddr = ""
		return
	}

	proxy.ShadowAddr = serviceAddr
}
