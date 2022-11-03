package transform

import "sync"

var Mutex = new(sync.Mutex)

var LocalAddrToRemoteAddr = map[string]string{}

func GetRemoteAddrFromLocalAddr(laddr string) (string, bool) {

	Mutex.Lock()
	defer Mutex.Unlock()
	raddr, ok := LocalAddrToRemoteAddr[laddr]
	return raddr, ok

}

func SetRemoteAddrToLocalAddr(laddr string, raddr string) {

	Mutex.Lock()
	defer Mutex.Unlock()
	LocalAddrToRemoteAddr[laddr] = raddr

}

func DeleteAddr(addr string) {
	Mutex.Lock()
	defer Mutex.Unlock()
	delete(LocalAddrToRemoteAddr, addr)
}
