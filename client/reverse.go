package client

import (
	"shadowproxy/config"
	"shadowproxy/logger"
	"shadowproxy/proxy"
	"strings"
)

func ReverseProxyClientRun() {

	client := proxy.RevProxyClient{
		ServerAddr: config.ShadowProxyConfig.ReverseServer,
		LinkAddr:   config.ShadowProxyConfig.ReverseLinkServer,
	}

	for _, v := range config.ShadowProxyConfig.ReverseRule {
		addrs := strings.Split(v, "->")
		logger.Log("REV CLI Start", v)
		go client.Link(addrs[0], addrs[1])

	}

}
