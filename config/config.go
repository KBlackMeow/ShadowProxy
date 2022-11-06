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
	EnableFilter  bool     `yaml:"enablefilter"`
	ConsoleOutput bool     `yaml:"consoleoutput"`
	Debug         bool     `yaml:"debug"`
	Client        bool     `yaml:"client"`
	Shadow        string   `yaml:"shadows"`
	AuthServer    string   `yaml:"authserver"`
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
		EnableFilter:  true,
		ConsoleOutput: true,
		Client:        false,
		Shadow:        "127.0.0.1:57575",
		AuthServer:    "127.0.0.1:5555",
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
