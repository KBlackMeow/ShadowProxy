package filter

import (
	"shadowproxy/config"
	"shadowproxy/logger"
	"strings"
	"sync"
	"time"
)

var EnableFillter bool = true

type IPStatue struct {
	IP        string
	Statu     byte
	BeginTime time.Time
	TTL       int64
}

func Filter(addr string) bool {

	addr = strings.Split(addr, ":")[0]
	// logger.Log("Check ", addr)
	var ret = false
	ret = ret || WhiteListFilter(addr)
	ret = ret || BlackListFilter(addr)
	ret = ret && config.ShadowProxyConfig.EnableFilter
	return ret

}

var IPStatuList = map[string]*IPStatue{}
var Mutex = new(sync.Mutex)

func AppendWhiteList(addr string) {

	addr = strings.Split(addr, ":")[0]

	Mutex.Lock()
	defer Mutex.Unlock()

	IP, ok := IPStatuList[addr]
	if ok {
		IP.Statu |= 1
		IP.BeginTime = time.Now()
		return
	}

	IPStatuList[addr] = &IPStatue{IP: addr, Statu: 1, BeginTime: time.Now(), TTL: 10000}

}

func WhiteListFilter(addr string) bool {

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

	addr = strings.Split(addr, ":")[0]

	Mutex.Lock()
	defer Mutex.Unlock()

	logger.Warn("Black list", addr, "appended")

	IP, ok := IPStatuList[addr]
	if ok {
		IP.Statu |= 2
		IP.BeginTime = time.Now()
		return
	}
	IPStatuList[addr] = &IPStatue{IP: addr, Statu: 2, BeginTime: time.Now(), TTL: 300000}

}

func BlackListFilter(addr string) bool {

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

func IPStatuLisClear() {

	for {
		for k, IP := range IPStatuList {
			if time.Since(IP.BeginTime).Milliseconds() > IP.TTL {
				logger.Log("Delete WhiteList IP: ", IP.IP)
				delete(IPStatuList, k)
			}
		}
		time.Sleep(time.Duration(500) * time.Millisecond)
	}

}

func InitFilter() {
	for _, v := range config.ShadowProxyConfig.WhiteList {
		AppendWhiteList(v)
	}
	for _, v := range config.ShadowProxyConfig.BlackList {
		AppendBlackList(v)
	}

	go IPStatuLisClear()
}
