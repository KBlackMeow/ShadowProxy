package service

import (
	"shadowproxy/proxy"
)

type RevProxyService struct {
	Service
	proxy.RevProxyServer
}

func (service RevProxyService) Run() {
	go service.RevProxyServer.Run()
}

func (service RevProxyService) GetAddr() string {
	return service.serviceAddr
}

func (service RevProxyService) GetName() string {
	return service.serviceName
}

func init() {
	service := RevProxyService{
		Service{serviceName: "reverse", serviceAddr: "0.0.0.0:50000"},
		proxy.RevProxyServer{ServerAddr: "0.0.0.0:50000", LinkAddr: "0.0.0.0:50001"},
	}
	ServiceAppend("reverse", service)
}
