package service

import (
	"encoding/json"
	"net"
	"shadowproxy/config"
	"shadowproxy/cryptotools"
	"shadowproxy/filter"
	"shadowproxy/ids"
	"shadowproxy/logger"
	"strconv"
	"strings"
	"time"
)

type AuthMessage struct {
	Token   string
	Pubkey  string
	Message string
}

type AuthService1 struct {
	Service
	listener *net.UDPConn
}

func (service AuthService1) token(remoteAddr string) string {

	return cryptotools.Hash_SHA512(remoteAddr)
}

func (service AuthService1) verifyToken(remoteAddr string, token string) bool {

	return token == cryptotools.Hash_SHA512(remoteAddr)

}

func (service AuthService1) auth() {

	for {
		buffer := make([]byte, 4096)
		n1, addr, err := service.listener.ReadFromUDP(buffer)
		ids.CheckIP(addr.String())
		if err != nil {
			logger.Error("UDP", err)
			return
		}

		loginMessage := AuthMessage{}
		e1 := json.Unmarshal(buffer[0:n1], &loginMessage)
		if e1 != nil {
			logger.Error(e1)
			continue
		}

		if len(loginMessage.Message) == 7 {

			Message := AuthMessage{Token: service.token(addr.String()), Pubkey: cryptotools.GetKey("public.pem")}
			data, _ := json.Marshal(&Message)
			go service.listener.WriteToUDP(data, addr)

		} else {

			cryptedMessage := loginMessage.Message
			Message := cryptotools.DecryptRSAToString(cryptedMessage)
			Messages := strings.Split(Message, "#")
			remoteAddr := addr.String()

			if Message == "" || len(Messages) != 3 {
				logger.Warn("Auth1", remoteAddr, "RSA Public Key is wrong")
				time.Sleep(time.Duration(3000) * time.Millisecond)
				continue
			}

			password := cryptotools.Hash_SHA512(Messages[0])
			MessageUnixTime, _ := strconv.ParseInt(Messages[1], 10, 64)
			MessageUnixTime = int64(MessageUnixTime)

			token := Messages[2]
			if !service.verifyToken(remoteAddr, token) {
				logger.Warn("Auth1", remoteAddr, "Token is wrong")
				time.Sleep(time.Duration(3000) * time.Millisecond)
				continue
			}

			if (time.Now().UnixMilli()-MessageUnixTime) > 0 && (time.Now().UnixMilli()-MessageUnixTime) < 1000 &&
				password == cryptotools.Hash_SHA512(config.ShadowProxyConfig.Password) {
				filter.AppendWhiteList(remoteAddr, 10000)
				continue
			}

			if password != cryptotools.Hash_SHA512(config.ShadowProxyConfig.Password) {
				logger.Warn("Auth1", remoteAddr, "Password is wrong")
			} else if (time.Now().UnixMilli() - MessageUnixTime) > 1000 {
				logger.Warn("Auth1", remoteAddr, "Unix Time exceed the time limit")
			} else {
				logger.Warn("Auth1", remoteAddr, "Alice is attacking the server")
			}

			time.Sleep(time.Duration(3000) * time.Millisecond)

		}

	}

}

func (service AuthService1) Run() {
	logger.Log("Auth1 Service Addr", service.serviceAddr)
	udpLAddr, _ := net.ResolveUDPAddr("udp4", service.serviceAddr)
	listener, err := net.ListenUDP("udp4", udpLAddr)
	if err != nil {
		logger.Error("UDP", err)
		return
	}
	service.listener = listener
	defer listener.Close()
	service.auth()

}

func (service AuthService1) GetName() string {

	return service.serviceName

}

func (service AuthService1) GetAddr() string {

	return service.serviceAddr

}

func init() {

	service := AuthService1{Service{serviceName: "auth1", serviceAddr: "0.0.0.0:5555"}, &net.UDPConn{}}
	ServiceAppend("auth1", service)

}
