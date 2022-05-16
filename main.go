package main

import (
	"flag"

	"shadowproxy/config"
	"shadowproxy/proxy"
	"shadowproxy/service"
)

func init() {
	service.StartServices()
	config.GetConfig()
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
		config.GetConfig()

	}

	config.InitComponentConfig()
	proxy.RunProxy()

}
