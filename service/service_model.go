package service

import (
	"shadowproxy/config"
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

var Services = map[string]Runner{}

func ServiceAppend(serviceName string, work Runner) {

	Services[serviceName] = work

}

func GetService(serviceName string) (Runner, bool) {
	service, ok := Services[serviceName]
	return service, ok
}

func InitServices() {

	for _, v := range config.ShadowProxyConfig.Services {
		service, ok := GetService(v)
		if ok {
			go service.Run()
		}
	}

}
