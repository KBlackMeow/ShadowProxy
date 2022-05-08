package main

import (
	"flag"
	"fmt"

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

	config.InitComponentConfig()
	proxy.RunProxy()

}
