package service

import (
	"fmt"
	"net/http"
	"shadowproxy/logger"
)

type FlagService struct {
	Service
}

func (service FlagService) flag(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is flag{W_W $_$}")
}

func (service FlagService) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/flag", service.flag)
	err := http.ListenAndServe(service.serviceAddr, mux)
	if err != nil {
		logger.Error(err)
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
	ServiceAppend(service)
}
