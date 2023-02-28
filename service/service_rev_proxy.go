package service

import (
	"shadowproxy/logger"
	"shadowproxy/proxy"
)

type RevProxyService struct {
	Service
	proxy.RevProxyServer
}

func (service RevProxyService) Run() {
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
		Service{serviceName: "reverse", serviceAddr: "0.0.0.0:20000"},
		proxy.RevProxyServer{ServerAddr: "0.0.0.0:20000", LinkAddr: "0.0.0.0:20001"},
	}
	ServiceAppend("reverse", service)
}
