package service

import "net/http"

type StaticService struct {
	Service
}

func (service StaticService) Contraller() {

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

}

func init() {

	service := StaticService{Service{serviceName: "http", serviceAddr: "127.0.0.1:57575"}}
	service.Contraller()

}
