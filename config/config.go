package config

import (
	"io/ioutil"
	"shadowproxy/cryptotools"
	"shadowproxy/fillter"
	"shadowproxy/logger"
	"shadowproxy/proxy"
	"shadowproxy/shadowtools"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BindAddr      string `yaml:"bindaddr"`
	BackendAddr   string `yaml:"backendaddr"`
	Protocol      string `yaml:"protocol"`
	Shadow        string `yaml:"shadow"`
	LogLevel      string `yaml:"loglevel"`
	Password      string `yaml:"password"`
	EnableFillter bool   `yaml:"enablefillter"`
	ConsoleOutput bool   `yaml:consoleoutput`
}

var FilePath = "config.yaml"
var ShadowProxyConfig Config

func GetConfig() {
	config := Config{}
	content, err := ioutil.ReadFile(FilePath)
	if err != nil {

		GenEmptyConfig()
		GetConfig()
		panic(err)
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
	ShadowProxyConfig = config
}

func GenEmptyConfig() {
	config := Config{
		BindAddr:      "0.0.0.0:30000",
		BackendAddr:   "127.0.0.1:40000",
		Protocol:      "tcp/udp",
		Shadow:        "auth",
		LogLevel:      "0",
		Password:      cryptotools.Md5_32("admin"),
		EnableFillter: true,
		ConsoleOutput: true,
	}
	content, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(FilePath, content, 0666)
}

func InitComponentConfig() {
	logger.ConsoleOutput = ShadowProxyConfig.ConsoleOutput
	fillter.EnableFillter = ShadowProxyConfig.EnableFillter

	shadowtools.SetShadowService(ShadowProxyConfig.Shadow)

	if num, err := strconv.ParseInt(ShadowProxyConfig.LogLevel, 10, 32); err == nil {
		if int(num) == 0 || int(num) == 1 || int(num) == 2 {
			logger.LogLevel = int(num)
		} else {
			logger.Error("LogLevel mast be 0, 1 or 2")
		}
	}

	proxy.ProxyProtocol = ShadowProxyConfig.Protocol
	proxy.ProxyBindAddr = ShadowProxyConfig.BindAddr
	proxy.ProxyBackendAddr = ShadowProxyConfig.BackendAddr

}
