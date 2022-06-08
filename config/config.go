package config

import (
	"io/ioutil"
	"shadowproxy/cryptotools"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogLevel      int      `yaml:"loglevel"`
	Password      string   `yaml:"password"`
	AuthSSL       bool     `yaml:"authssl"`
	EnableFillter bool     `yaml:"enablefillter"`
	ConsoleOutput bool     `yaml:"consoleoutput"`
	Debug         bool     `yaml:"debug"`
	Shadow        string   `yaml:"shadows"`
	Services      []string `yaml:"services"`
	Rules         []string `yaml:"rules"`
	WhiteList     []string `yaml:"whitelist"`
	BlackList     []string `yaml:"blacklist"`
	CMD           []string `yaml:"cmd"`
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

		LogLevel:      0,
		Password:      cryptotools.Hash_MD5("admin"),
		AuthSSL:       false,
		EnableFillter: true,
		ConsoleOutput: true,
		Shadow:        "auth",
		Services:      []string{"auth", "flag", "cmd"},
		Rules:         []string{"tcp://0.0.0.0:30000->127.0.0.1:40000"},
		WhiteList:     []string{"127.0.0.1"},
		BlackList:     []string{},
		CMD:           []string{},
	}
	content, err := yaml.Marshal(ShadowProxyConfig)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(FilePath, content, 0666)

}
