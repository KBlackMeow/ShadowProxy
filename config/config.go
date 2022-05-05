package config

import (
	"io/ioutil"
	"shadowproxy/cryptotools"
	"shadowproxy/logger"

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
}

var FilePath = "config.yaml"
var ShadowProxyConfig Config

func GetConfig() {
	config := Config{}
	content, err := ioutil.ReadFile(FilePath)
	if err != nil {
		logger.Error(err)
		GenEmptyConfig()
		GetConfig()
		return
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		logger.Error(err)
		return
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
		EnableFillter: true}
	content, err := yaml.Marshal(config)
	if err != nil {
		logger.Error(err)
	}
	ioutil.WriteFile(FilePath, content, 0666)
	logger.Log("config.yaml has been created")
}
