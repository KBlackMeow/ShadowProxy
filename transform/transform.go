package transform

import "sync"

var Mutex = new(sync.Mutex)

var LAddrToRAddr = map[string]string{}

func GetRAddrFromLAddr(laddr string) (string, bool) {

	Mutex.Lock()
	defer Mutex.Unlock()
	raddr, ok := LAddrToRAddr[laddr]
	return raddr, ok

}

func SetRAddrToLAddr(laddr string, raddr string) {

	Mutex.Lock()
	defer Mutex.Unlock()
	LAddrToRAddr[laddr] = raddr

}

func DeleteAddr(addr string) {
	Mutex.Lock()
	defer Mutex.Unlock()
	delete(LAddrToRAddr, addr)
}

var NameToAddr = map[string]string{}
