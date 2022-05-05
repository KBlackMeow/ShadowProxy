package main

import (
	"flag"
	"fmt"
	"strconv"

	"shadowproxy/config"
	"shadowproxy/fillter"
	"shadowproxy/logger"
	"shadowproxy/proxy"
	"shadowproxy/service"
	"shadowproxy/shadowtools"
)

func init() {
	service.StartServices()
	config.GetConfig()
}
func main() {

	help := flag.Bool("help", false, "print usage")
	cfg := flag.String("config", "", "use config file")
	bind := flag.String("bind", "127.0.0.1:30000", "The address to bind to")
	backend := flag.String("backend", "", "The backend server address")
	protocol := flag.String("protocol", "tcp", "To use tcp, udp or tcp/udp")
	shadow := flag.String("shadow", "", "The shadow server address")
	loglevel := flag.Int("loglevel", 0, "The level of log 0:all, 1 warn and error,2 error")
	password := flag.String("password", "admin", "password of auth system,and it must not be empty")
	enableFillter := flag.Bool("fillter", true, "enable fillter")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *cfg != "" {

		config.FilePath = *cfg
		config.GetConfig()
		config.ShadowProxyConfig.BindAddr = *bind
		config.ShadowProxyConfig.BackendAddr = *backend
		config.ShadowProxyConfig.Protocol = *protocol
		config.ShadowProxyConfig.Shadow = *shadow
		config.ShadowProxyConfig.LogLevel = fmt.Sprint(*loglevel)
		config.ShadowProxyConfig.Password = *password
		config.ShadowProxyConfig.EnableFillter = *enableFillter
	}

	fillter.EnableFillter = config.ShadowProxyConfig.EnableFillter

	num, err := strconv.ParseInt(config.ShadowProxyConfig.LogLevel, 10, 32)
	if err == nil {
		if int(num) == 0 || int(num) == 1 || int(num) == 2 {
			logger.LogLevel = int(num)
		} else {
			logger.Error("LogLevel mast be 0, 1 or 2")
		}
	}

	shadowtools.SetShadowService(config.ShadowProxyConfig.Shadow)

	if config.ShadowProxyConfig.Protocol == "tcp" {
		proxy.WG.Add(1)
		go proxy.RunTPortProxy(config.ShadowProxyConfig.BindAddr, config.ShadowProxyConfig.BackendAddr)

	} else if config.ShadowProxyConfig.Protocol == "udp" {
		proxy.WG.Add(1)
		go proxy.RunUPortProxy(config.ShadowProxyConfig.BindAddr, config.ShadowProxyConfig.BackendAddr)

	} else if config.ShadowProxyConfig.Protocol == "tcp/udp" {
		proxy.WG.Add(2)
		go proxy.RunTPortProxy(config.ShadowProxyConfig.BindAddr, config.ShadowProxyConfig.BackendAddr)
		go proxy.RunUPortProxy(config.ShadowProxyConfig.BindAddr, config.ShadowProxyConfig.BackendAddr)

	} else {
		flag.Usage()
		return
	}

	proxy.WG.Wait()
}
