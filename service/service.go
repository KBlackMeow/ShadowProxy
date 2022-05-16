package service

import (
	"shadowproxy/proxy"
)

type Service struct {
	serviceName string
	serviceAddr string
}

type Runner interface {
	Run()
	GetName() string
	GetAddr() string
}

var Services []Runner

func ServiceAppend(work Runner) {
	Services = append(Services, work)
}

func InitServices() {
	for _, service := range Services {
		go service.Run()
		proxy.NameToAddr[service.GetName()] = service.GetAddr()
	}
}
