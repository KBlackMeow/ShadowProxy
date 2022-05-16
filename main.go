package main

import (
	"flag"

	"shadowproxy/config"
	"shadowproxy/proxy"
	"shadowproxy/service"
	"shadowproxy/shadowtools"
)

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
	config.InitConfig()
	service.InitServices()
	shadowtools.InitShadowService()
	proxy.RunProxy()

}
