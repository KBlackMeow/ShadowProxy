package fillter

import (
	"shadowproxy/logger"
	"strings"
)

var EnableFillter bool = true

type IPStatue struct {
	IP    string
	Statu byte
}

func Fillter(addr string) bool {

	if tem := strings.Split(addr, ":"); len(tem) == 2 {
		addr = tem[0]
	}

	var ret = false
	ret = ret || WhiteListFillter(addr)
	ret = ret || BlackListFillter(addr)
	ret = ret && EnableFillter
	return ret
}

var IPStatuList = map[string]*IPStatue{}

func AppendWhiteList(addr string) {
	if tem := strings.Split(addr, ":"); len(tem) == 2 {
		addr = tem[0]
	}

	IP, ok := IPStatuList[addr]
	if ok {
		IP.Statu |= 1
		return
	}

	IPStatuList[addr] = &IPStatue{IP: addr, Statu: 1}
}

func WhiteListFillter(addr string) bool {

	IP, ok := IPStatuList[addr]
	if ok {
		if IP.Statu%2 == 1 {
			return false
		}
	}
	return true
}

func AppendBlackList(addr string) {
	if tem := strings.Split(addr, ":"); len(tem) == 2 {
		addr = tem[0]
	}

	logger.Warn("Black list", addr, "appended")

	IP, ok := IPStatuList[addr]
	if ok {
		IP.Statu |= 2
		return
	}
	IPStatuList[addr] = &IPStatue{IP: addr, Statu: 2}

}

func BlackListFillter(addr string) bool {

	IP, ok := IPStatuList[addr]
	if ok {
		if (IP.Statu%4)/2 == 1 {
			return true
		}
	}

	return false
}
