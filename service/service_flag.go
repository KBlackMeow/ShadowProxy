package service

import (
	"fmt"
	"net/http"
	"shadowproxy/config"
	"shadowproxy/cryptotools"
	"shadowproxy/logger"
)

type FlagService struct {
	Service
}

func (service FlagService) flag(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "flag{"+cryptotools.Hash_MD5("flag")+"}")

}

func (service FlagService) Run() {

	logger.Log("Flag Service Addr", service.serviceAddr)
	mux := http.NewServeMux()
	mux.HandleFunc("/flag", service.flag)

	if config.ShadowProxyConfig.AuthSSL {
		err := http.ListenAndServeTLS(service.serviceAddr, "server.crt", "server.key", mux)
		if err != nil {
			logger.Error(err)
		}
	} else {
		err := http.ListenAndServe(service.serviceAddr, mux)
		if err != nil {
			logger.Error(err)
		}
	}

}

func (service FlagService) GetAddr() string {

	return service.serviceAddr

}

func (service FlagService) GetName() string {

	return service.serviceName

}

func init() {

	service := FlagService{Service{serviceAddr: "127.0.0.1:40000", serviceName: "flag"}}
	ServiceAppend("flag", service)

}
