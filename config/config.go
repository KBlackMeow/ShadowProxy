package config

import (
	"io/ioutil"
	"shadowproxy/cryptotools"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BindAddr      string `yaml:"bindaddr"`
	BackendAddr   string `yaml:"backendaddr"`
	Protocol      string `yaml:"protocol"`
	Shadow        string `yaml:"shadow"`
	LogLevel      int    `yaml:"loglevel"`
	Password      string `yaml:"password"`
	EnableFillter bool   `yaml:"enablefillter"`
	ConsoleOutput bool   `yaml:"consoleoutput"`
	Debug         bool   `yaml:"debug"`
}

var FilePath = "config.yaml"
var ShadowProxyConfig Config

func InitConfig() {

	config := Config{}
	content, err := ioutil.ReadFile(FilePath)
	if err != nil {
		GenEmptyConfig()
		return
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
	ShadowProxyConfig = config

}

func GenEmptyConfig() {

	ShadowProxyConfig = Config{

		BindAddr:      "0.0.0.0:30000",
		BackendAddr:   "127.0.0.1:40000",
		Protocol:      "tcp/udp",
		Shadow:        "auth",
		LogLevel:      0,
		Password:      cryptotools.Hash_MD5("admin"),
		EnableFillter: true,
		ConsoleOutput: true,
	}
	content, err := yaml.Marshal(ShadowProxyConfig)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(FilePath, content, 0666)

}
