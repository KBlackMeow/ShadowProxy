package config

type TempConfig struct {
	Key    string
	AES_IV string
	PubKey string
	PriKey string
}

var TempCfgObj = TempConfig{
	Key:    "12345678901234567890123456789012",
	AES_IV: "4324frgt5tg534rd",
}
