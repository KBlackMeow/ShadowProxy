package service

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"shadowproxy/config"
	"shadowproxy/cryptotools"
	"shadowproxy/filter"
	"shadowproxy/logger"
	"strconv"
	"strings"
	"time"
)

type LoginInfo struct {
	CryptedMessage string `json:"CryptedMessage"`
}

type UserInfo struct {
	UserAddr      string
	UserLoginTime string
}

type AuthService2 struct {
	Service
}

func (service AuthService2) token(remoteAddr string) string {

	return cryptotools.Hash_SHA512(remoteAddr)
}

func (service AuthService2) verifyToken(remoteAddr string, token string) bool {

	return token == cryptotools.Hash_SHA512(remoteAddr)

}

func (service AuthService2) verify(w http.ResponseWriter, r *http.Request) {

	remoteAddr := r.RemoteAddr

	var loginfo LoginInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginfo)

	if err != nil {
		logger.Error(err)
		time.Sleep(time.Duration(3000) * time.Millisecond)
		return
	}

	CryptedMessage := loginfo.CryptedMessage
	Message := cryptotools.DecryptRSAToString(CryptedMessage)
	Messages := strings.Split(Message, "#")

	if Message == "" || len(Messages) != 3 {
		logger.Warn("Auth2", remoteAddr, "RSA Public Key is wrong")
		time.Sleep(time.Duration(3000) * time.Millisecond)
		return
	}

	password := cryptotools.Hash_SHA512(Messages[0])
	MessageUnixTime, _ := strconv.ParseInt(Messages[1], 10, 64)
	MessageUnixTime = int64(MessageUnixTime)

	token := Messages[2]
	if !service.verifyToken(remoteAddr, token) {
		logger.Warn("Auth2", remoteAddr, "Token is wrong")
		time.Sleep(time.Duration(3000) * time.Millisecond)
		return
	}

	if (time.Now().UnixMilli()-MessageUnixTime) > 0 && (time.Now().UnixMilli()-MessageUnixTime) < 1000 &&
		password == cryptotools.Hash_SHA512(config.ShadowProxyConfig.Password) {
		filter.AppendWhiteList(remoteAddr, 10000)

		userinfo := UserInfo{UserAddr: remoteAddr, UserLoginTime: logger.TimeNow()}
		res, _ := json.Marshal(&userinfo)
		fmt.Fprint(w, string(res))
		return
	}

	if password != cryptotools.Hash_SHA512(config.ShadowProxyConfig.Password) {
		logger.Warn("Auth2", remoteAddr, "Password is wrong")
	} else if (time.Now().UnixMilli() - MessageUnixTime) > 1000 {
		logger.Warn("Auth2", remoteAddr, "Unix Time exceed the time limit")
	} else {
		logger.Warn("Auth2", remoteAddr, "Alice is attacking the server")
	}

	time.Sleep(time.Duration(3000) * time.Millisecond)
	userinfo := UserInfo{}
	data, _ := json.Marshal(&userinfo)
	fmt.Fprint(w, string(data))

}

func (service AuthService2) auth(w http.ResponseWriter, r *http.Request) {

	temp, err := template.ParseFiles("template/auth2.html")
	if err != nil {
		logger.Error(err)
		return
	}

	type TempleInfo struct {
		PubKey string
		Token  string
	}
	remoteAddr := r.RemoteAddr

	x := TempleInfo{PubKey: cryptotools.GetKey("public.pem"), Token: service.token(remoteAddr)}
	temp.Execute(w, x)

}

func (service AuthService2) Contraller() {

	http.HandleFunc("/auth2", service.auth)
	http.HandleFunc("/verify2", service.verify)

}

func (service AuthService2) Run() {

	if config.ShadowProxyConfig.AuthServer != "" {
		service.serviceAddr = config.ShadowProxyConfig.AuthServer
	}

	logger.Log("Auth2 Service Addr", service.serviceAddr)
	if config.ShadowProxyConfig.AuthSSL {
		err := http.ListenAndServeTLS(service.serviceAddr, "server.crt", "server.key", nil)
		if err != nil {
			logger.Error(err)
		}
	} else {
		err := http.ListenAndServe(service.serviceAddr, nil)
		if err != nil {
			logger.Error(err)
		}
	}

}

func (service AuthService2) GetName() string {

	return service.serviceName

}

func (service AuthService2) GetAddr() string {

	return service.serviceAddr

}

func init() {

	service := AuthService2{Service{serviceName: "auth2", serviceAddr: "127.0.0.1:50002"}}
	service.Contraller()
	ServiceAppend("auth2", service)

}
