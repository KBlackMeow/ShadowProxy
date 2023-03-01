package config

import (
	"os"
	"shadowproxy/cryptotools"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LogLevel          int      `yaml:"loglevel"`
	Password          string   `yaml:"password"`
	AuthSSL           bool     `yaml:"authssl"`
	EnableFilter      bool     `yaml:"enablefilter"`
	ConsoleOutput     bool     `yaml:"consoleoutput"`
	Debug             bool     `yaml:"debug"`
	Client            bool     `yaml:"client"`
	Shadow            string   `yaml:"shadows"`
	AuthServer        string   `yaml:"authserver"`
	ReverseServer     string   `yaml:"revserver"`
	ReverseLinkServer string   `yaml:"reverselinkserver"`
	ReverseRule       []string `yaml:"reverserule"`
	Services          []string `yaml:"services"`
	Rules             []string `yaml:"rules"`
	WhiteList         []string `yaml:"whitelist"`
	BlackList         []string `yaml:"blacklist"`
	CMD               []string `yaml:"cmd"`
}

var FilePath = "config.yaml"
var ShadowProxyConfig Config

func InitConfig() {

	config := Config{}
	content, err := os.ReadFile(FilePath)
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

		LogLevel:          0,
		Password:          cryptotools.Hash_MD5("admin"),
		AuthSSL:           false,
		EnableFilter:      true,
		ConsoleOutput:     true,
		Client:            false,
		Shadow:            "127.0.0.1:57575",
		AuthServer:        "127.0.0.1:50002",
		ReverseServer:     "127.0.0.1:50000",
		ReverseLinkServer: "127.0.0.1:50001",
		ReverseRule:       []string{"127.0.0.1:41000->127.0.0.1:41001"},
		Services:          []string{"auth1", "auth2", "flag", "cmd", "reverse"},
		Rules:             []string{"tcp://127.0.0.1:30000->127.0.0.1:40000"},
		WhiteList:         []string{"127.0.0.1"},
		BlackList:         []string{},
		CMD:               []string{},
	}
	content, err := yaml.Marshal(ShadowProxyConfig)
	if err != nil {
		panic(err)
	}
	os.WriteFile(FilePath, content, 0666)

}
