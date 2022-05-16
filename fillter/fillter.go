package fillter

import (
	"net"
	"shadowproxy/config"
	"shadowproxy/logger"
	"sync"
)

var EnableFillter bool = true

type IPStatue struct {
	IP    string
	Statu byte
}

func Fillter(addr string) bool {

	addr = net.ParseIP(addr).String()

	var ret = false
	ret = ret || WhiteListFillter(addr)
	ret = ret || BlackListFillter(addr)
	ret = ret && config.ShadowProxyConfig.EnableFillter
	return ret
}

var IPStatuList = map[string]*IPStatue{}
var Mutex = new(sync.Mutex)

func AppendWhiteList(addr string) {

	addr = net.ParseIP(addr).String()

	Mutex.Lock()
	defer Mutex.Unlock()

	IP, ok := IPStatuList[addr]
	if ok {
		IP.Statu |= 1
		return
	}

	IPStatuList[addr] = &IPStatue{IP: addr, Statu: 1}
}

func WhiteListFillter(addr string) bool {

	Mutex.Lock()
	defer Mutex.Unlock()

	IP, ok := IPStatuList[addr]
	if ok {
		if IP.Statu%2 == 1 {
			return false
		}
	}
	return true
}

func AppendBlackList(addr string) {

	addr = net.ParseIP(addr).String()

	Mutex.Lock()
	defer Mutex.Unlock()

	logger.Warn("Black list", addr, "appended")

	IP, ok := IPStatuList[addr]
	if ok {
		IP.Statu |= 2
		return
	}
	IPStatuList[addr] = &IPStatue{IP: addr, Statu: 2}

}

func BlackListFillter(addr string) bool {

	Mutex.Lock()
	defer Mutex.Unlock()

	IP, ok := IPStatuList[addr]
	if ok {
		if (IP.Statu%4)/2 == 1 {
			return true
		}
	}

	return false
}
