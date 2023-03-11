package config

type TempConfig struct {
	Key    string
	PubKey string
	PriKey string
}

var TempCfgObj = TempConfig{
	Key: "12345678901234567890123456789012",
}
