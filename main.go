package main

import (
	"flag"

	"shadowproxy/config"
	"shadowproxy/filter"
	"shadowproxy/proxy"
	"shadowproxy/service"
	"shadowproxy/shadowtools"
)

func init() {
	config.InitConfig()
}

func ClientComponentInit() {
	// client.ClientInit()
}

func ServerComponentInit() {

	service.InitServices()
	shadowtools.InitShadowService()
	filter.InitFilter()
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
