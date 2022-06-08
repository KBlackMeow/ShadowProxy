package main

import (
	"flag"

	"shadowproxy/config"
	"shadowproxy/fillter"
	"shadowproxy/proxy"
	"shadowproxy/service"
	"shadowproxy/shadowtools"
	"shadowproxy/tunnel"
)

func init() {
	config.InitConfig()
	tunnel.Run()
}

func ComponentInit() {

	service.InitServices()
	shadowtools.InitShadowService()
	fillter.InitFillter()
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

	ComponentInit()
}
