package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"shadowproxy/config"
	"shadowproxy/cryptotools"
	"shadowproxy/logger"
	"time"
)

type AuthMessage struct {
	Token   string
	Pubkey  string
	Message string
}
type Client struct {
	Addr     string
	Token    string
	Password string
	Conn     net.Conn
}

func (c *Client) GetKey() {

	Message := AuthMessage{
		Token:   "",
		Pubkey:  "",
		Message: "1234567",
	}

	data, e1 := json.Marshal(Message)
	if e1 != nil {
		logger.Error(e1)

	}
	c.Conn.Write(data)

}

func (c *Client) Listen() {

	for {
		buffer := make([]byte, 4096)
		n1, err := c.Conn.Read(buffer)
		if err != nil {
			logger.Error(err)
			continue
		}

		Message := AuthMessage{}

		e1 := json.Unmarshal(buffer[:n1], &Message)

		if e1 != nil {
			logger.Error(e1)
			continue
		}

		if Message.Pubkey != "" {
			c.Token = Message.Token
			config.TempCfgObj.PubKey = Message.Pubkey
			logger.Log("Login : Get PubKey, length:", len(config.TempCfgObj.PubKey))
		}

	}

}

func (c Client) Login() {
	conn, err := net.Dial("udp", c.Addr)
	if err != nil {
		logger.Error(err)
		return
	}
	c.Conn = conn
	defer c.Conn.Close()
	go c.Listen()

	for {
		if config.TempCfgObj.PubKey != "" {
			Message := c.Password + "#" + fmt.Sprint(time.Now().UnixMilli()) + "#" + c.Token

			cipherText := cryptotools.RSA_Encode([]byte(Message), config.TempCfgObj.PubKey)
			cMessageb64 := base64.StdEncoding.EncodeToString(cipherText)

			loginMessage := AuthMessage{}
			loginMessage.Message = cMessageb64
			data, _ := json.Marshal(loginMessage)

			c.Conn.Write(data)

		} else {
			c.GetKey()
		}

		time.Sleep(time.Duration(3000) * time.Millisecond)

	}
}

func ClientRun() {

	c := Client{Token: "", Password: config.ShadowProxyConfig.Password, Addr: config.ShadowProxyConfig.AuthServer}
	config.TempCfgObj.Key = cryptotools.Hash_MD5("4455667")
	c.Login()
}
