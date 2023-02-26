package main

import (
	"flag"
	"time"

	"shadowproxy/config"
	"shadowproxy/proxy"
	"shadowproxy/tunnel"
)

func init() {
	config.InitConfig()
}

func ClientComponentInit() {
	// client.ClientInit()
}

func ServerComponentInit() {

	// service.InitServices()
	// shadowtools.InitShadowService()
	// filter.InitFilter()

	tunnel.TunnelInit2()
	time.Sleep(time.Second * 1)
	tunnel.TunnelInit1()

	proxy.RunProxy()

}

func main() {
	help := flag.Bool("help", false, "print usage")
	cfg := flag.String("config", "", "use config file")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *cfg != "" {
		config.FilePath = *cfg
	}

	if !config.ShadowProxyConfig.Client {
		ServerComponentInit()
	} else {
		ClientComponentInit()
	}
}
