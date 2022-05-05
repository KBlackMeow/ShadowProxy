package fillter

import (
	"shadowproxy/logger"
	"strings"
)

var EnableFillter bool = true

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

var WhiteList = []string{}

func AppendWhiteList(addr string) {
	if tem := strings.Split(addr, ":"); len(tem) == 2 {
		addr = tem[0]
	}
	WhiteList = append(WhiteList, addr)
}

func WhiteListFillter(addr string) bool {

	for _, name := range WhiteList {
		if addr == name {
			return false
		}
	}
	return true
}

var BlackList = []string{}

func AppendBlackList(addr string) {
	if tem := strings.Split(addr, ":"); len(tem) == 2 {
		addr = tem[0]
	}
	BlackList = append(BlackList, addr)
	logger.Warn("Black list", addr, "appended")
}

func BlackListFillter(addr string) bool {

	for _, name := range BlackList {
		if addr == name {
			return true
		}
	}
	return false
}
