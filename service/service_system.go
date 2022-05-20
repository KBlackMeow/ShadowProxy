package service

import (
	"os/exec"
	"shadowproxy/config"
	"shadowproxy/logger"
	"strings"
)

type SystemService struct {
	Service
}

func (service SystemService) Run() {

	for _, v := range config.ShadowProxyConfig.CMD {
		args := strings.Split(v, " ")
		command := exec.Command(args[0], args[1:]...)
		err := command.Start()
		if err != nil {
			logger.Error("cmd", err)
		}
	}

}

func (service SystemService) GetAddr() string {

	return service.serviceAddr

}

func (service SystemService) GetName() string {

	return service.serviceName

}

func init() {

	service := SystemService{Service{serviceName: "cmd", serviceAddr: ""}}
	ServiceAppend("cmd", service)

}
