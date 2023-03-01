package service

import (
	"shadowproxy/config"
	"shadowproxy/logger"
	"shadowproxy/proxy"
)

type RevProxyService struct {
	Service
	proxy.RevProxyServer
}

func (service RevProxyService) Run() {

	if config.ShadowProxyConfig.ReverseServer != "" {
		service.serviceAddr = config.ShadowProxyConfig.ReverseServer
		service.RevProxyServer.ServerAddr = config.ShadowProxyConfig.ReverseServer
	}

	if config.ShadowProxyConfig.ReverseLinkServer != "" {
		service.RevProxyServer.LinkAddr = config.ShadowProxyConfig.ReverseLinkServer
	}

	go service.RevProxyServer.Run()
	logger.Log("Reverse Service Start", service.serviceAddr)
}

func (service RevProxyService) GetAddr() string {
	return service.serviceAddr
}

func (service RevProxyService) GetName() string {
	return service.serviceName
}

func init() {
	service := RevProxyService{
		Service{serviceName: "reverse", serviceAddr: "127.0.0.1:50000"},
		proxy.RevProxyServer{ServerAddr: "127.0.0.1:50000", LinkAddr: "127.0.0.1:50001"},
	}
	ServiceAppend("reverse", service)
}
