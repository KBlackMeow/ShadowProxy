package client

import (
	"shadowproxy/proxy"
)

func ReverseProxyClientRun() {

	client := proxy.RevProxyClient{
		ServerAddr: "0.0.0.0:20000",
		LinkAddr:   "0.0.0.0:20001",
	}
	// time.Sleep(time.Second * 1)
	go client.Run()
}
