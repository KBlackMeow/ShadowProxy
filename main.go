package main

import (
	"flag"
	"sync"
	"time"

	"shadowproxy/client"
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
	client.ClientRun()

}

func ServerComponentInit() {
	service.InitServices()
	shadowtools.InitShadowService()
	filter.InitFilter()
	proxy.RunProxy()

}

func TEST() {
	// TEST
	time.Sleep(time.Second * 1)
	client.ReverseProxyClientRun()
	client.ClientRun()
}

func main() {
	help := flag.Bool("help", false, "print usage")
	cfg := flag.String("config", "", "use config file")
	test := flag.Bool("test", false, "print usage")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *cfg != "" {
		config.FilePath = *cfg
	}

	if !config.ShadowProxyConfig.Client {
		go ServerComponentInit()
	} else {
		go ClientComponentInit()
	}

	if *test {
		TEST()
	}

	var WG sync.WaitGroup
	WG.Add(1)
	WG.Wait()
}
