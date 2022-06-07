package main

import (
	"flag"

	"shadowproxy/config"
	"shadowproxy/fillter"
	"shadowproxy/proxy"
	"shadowproxy/service"
	"shadowproxy/shadowtools"
)

func init() {

}

func ComponentInit() {
	config.InitConfig()
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
