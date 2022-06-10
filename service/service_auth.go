package service

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"shadowproxy/config"
	"shadowproxy/cryptotools"
	"shadowproxy/fillter"
	"shadowproxy/logger"
	"shadowproxy/proxy"
	"shadowproxy/transform"
	"strconv"
	"strings"
	"time"
)

type AuthService struct {
	Service
}

type LoginInfo struct {
	CMsg string `json:"cmsg"`
}

type UserInfo struct {
	UserAddr      string
	UserLoginTime string
}

func (service AuthService) token(remoteAddr string) string {

	remoteAddr, ok := transform.GetRAddrFromLAddr(remoteAddr)
	if ok {
		return cryptotools.Hash_SHA512(remoteAddr)
	}
	return ""

}

func (service Service) verifyToken(remoteAddr string, token string) bool {

	if token == cryptotools.Hash_SHA512(remoteAddr) {
		return true
	}
	return false

}

func (service AuthService) verify(w http.ResponseWriter, r *http.Request) {

	remoteAddr, ok := transform.GetRAddrFromLAddr(r.RemoteAddr)

	if ok {
		var loginfo LoginInfo
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&loginfo)

		if err != nil {
			logger.Error(err)
			time.Sleep(time.Duration(3000) * time.Millisecond)
			return
		}

		cmsg := loginfo.CMsg
		msg := cryptotools.DecryptRSAToString(cmsg)
		msgs := strings.Split(msg, "#")

		if msg == "" || len(msgs) != 3 {
			logger.Warn("Auth", remoteAddr, "RSA Public Key is wrong")
			time.Sleep(time.Duration(3000) * time.Millisecond)
			return
		}

		password := cryptotools.Hash_SHA512(msgs[0])
		msgUnixTime, _ := strconv.ParseInt(msgs[1], 10, 64)
		msgUnixTime = int64(msgUnixTime)

		token := msgs[2]
		if !service.verifyToken(remoteAddr, token) {
			logger.Warn("Auth", remoteAddr, "Token is wrong")
			time.Sleep(time.Duration(3000) * time.Millisecond)
			return
		}

		if (time.Now().UnixMilli()-msgUnixTime) > 0 && (time.Now().UnixMilli()-msgUnixTime) < 1000 &&
			password == cryptotools.Hash_SHA512(config.ShadowProxyConfig.Password) {
			fillter.AppendWhiteList(remoteAddr)

			userinfo := UserInfo{UserAddr: remoteAddr, UserLoginTime: logger.TimeNow()}
			res, _ := json.Marshal(&userinfo)
			fmt.Fprintf(w, string(res))

			go proxy.ClearConnFromIP(remoteAddr, 2000)
			return
		}

		if password != cryptotools.Hash_SHA512(config.ShadowProxyConfig.Password) {
			logger.Warn("Auth", remoteAddr, "Password is wrong")
		} else if (time.Now().UnixMilli() - msgUnixTime) > 1000 {
			logger.Warn("Auth", remoteAddr, "Unix Time exceed the time limit")
		} else {
			logger.Warn("Auth", remoteAddr, "Alice is attacking the server")
		}
	}

	time.Sleep(time.Duration(3000) * time.Millisecond)
	userinfo := UserInfo{}
	data, _ := json.Marshal(&userinfo)
	fmt.Fprintf(w, string(data))

}

func (service AuthService) auth(w http.ResponseWriter, r *http.Request) {

	temp, err := template.ParseFiles("template/auth.html")
	if err != nil {
		logger.Error(err)
		return
	}

	type TempleInfo struct {
		PubKey string
		Token  string
	}
	x := TempleInfo{PubKey: cryptotools.GetKey("public.pem"), Token: service.token(r.RemoteAddr)}
	temp.Execute(w, x)

}

func (service AuthService) Contraller() {

	http.HandleFunc("/auth", service.auth)
	http.HandleFunc("/verify", service.verify)

}

func (service AuthService) Run() {

	logger.Log("Auth Service Starting...")
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

func (service AuthService) GetName() string {

	return service.serviceName

}

func (service AuthService) GetAddr() string {

	return service.serviceAddr

}

func init() {

	service := AuthService{Service{serviceName: "auth", serviceAddr: "127.0.0.1:57575"}}
	service.Contraller()
	ServiceAppend("auth", service)

}
