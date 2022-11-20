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

type AuthMsg struct {
	Token  string
	Pubkey string
	Msg    string
}

type AuthService3 struct {
	Service
	listener *net.UDPConn
}

func (service AuthService3) token(remoteAddr string) string {

	return cryptotools.Hash_SHA512(remoteAddr)
}

func (service AuthService3) verifyToken(remoteAddr string, token string) bool {

	return token == cryptotools.Hash_SHA512(remoteAddr)

}

func (service AuthService3) auth() {

	for {
		buffer := make([]byte, 4096)
		n1, addr, err := service.listener.ReadFromUDP(buffer)
		ids.CheckIP(addr.String())
		if err != nil {
			logger.Error("UDP", err)
			return
		}

		loginMsg := AuthMsg{}
		e1 := json.Unmarshal(buffer[0:n1], &loginMsg)
		if e1 != nil {
			logger.Error(e1)
			continue
		}

		if len(loginMsg.Msg) == 7 {

			msg := AuthMsg{Token: service.token(addr.String()), Pubkey: cryptotools.GetKey("public.pem")}
			data, _ := json.Marshal(&msg)
			go service.listener.WriteToUDP(data, addr)

		} else {

			cmsg := loginMsg.Msg
			msg := cryptotools.DecryptRSAToString(cmsg)
			msgs := strings.Split(msg, "#")
			remoteAddr := addr.String()

			if msg == "" || len(msgs) != 3 {
				logger.Warn("Auth3", remoteAddr, "RSA Public Key is wrong")
				time.Sleep(time.Duration(3000) * time.Millisecond)
				continue
			}

			password := cryptotools.Hash_SHA512(msgs[0])
			msgUnixTime, _ := strconv.ParseInt(msgs[1], 10, 64)
			msgUnixTime = int64(msgUnixTime)

			token := msgs[2]
			if !service.verifyToken(remoteAddr, token) {
				logger.Warn("Auth3", remoteAddr, "Token is wrong")
				time.Sleep(time.Duration(3000) * time.Millisecond)
				continue
			}

			if (time.Now().UnixMilli()-msgUnixTime) > 0 && (time.Now().UnixMilli()-msgUnixTime) < 1000 &&
				password == cryptotools.Hash_SHA512(config.ShadowProxyConfig.Password) {
				filter.AppendWhiteList(remoteAddr, 10000)
				continue
			}

			if password != cryptotools.Hash_SHA512(config.ShadowProxyConfig.Password) {
				logger.Warn("Auth3", remoteAddr, "Password is wrong")
			} else if (time.Now().UnixMilli() - msgUnixTime) > 1000 {
				logger.Warn("Auth3", remoteAddr, "Unix Time exceed the time limit")
			} else {
				logger.Warn("Auth3", remoteAddr, "Alice is attacking the server")
			}

			time.Sleep(time.Duration(3000) * time.Millisecond)

		}

	}

}

func (service AuthService3) Run() {
	logger.Log("Auth3 Service Addr", service.serviceAddr)
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

func (service AuthService3) GetName() string {

	return service.serviceName

}

func (service AuthService3) GetAddr() string {

	return service.serviceAddr

}

func init() {

	service := AuthService3{Service{serviceName: "auth3", serviceAddr: "0.0.0.0:5555"}, &net.UDPConn{}}
	ServiceAppend("auth3", service)

}
