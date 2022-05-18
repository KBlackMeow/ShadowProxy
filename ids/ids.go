package ids

import (
	"net"
	"shadowproxy/fillter"
	"sync"
	"time"
)

var GuestMap = map[string]*Guest{}
var Mutex = new(sync.Mutex)

type Guest struct {
	GuestIP              string
	VisitTimeRecorder    []int
	PackageLengthRecoder []int
}

func (guest *Guest) visited() {

	guest.VisitTimeRecorder = append(guest.VisitTimeRecorder, int(time.Now().UnixMilli()))

}

func (guest *Guest) sent(length int) {

	guest.PackageLengthRecoder = append(guest.PackageLengthRecoder, length)

}

func (guest Guest) count() int {

	ret := 0
	nowTime := int(time.Now().UnixMilli())

	for _, v := range guest.VisitTimeRecorder {

		if nowTime-v < 500 {
			ret++
		}
	}

	length := len(guest.VisitTimeRecorder)
	var index int
	if length < 20 {
		index = 0
	} else {
		index = length - 20
	}
	guest.VisitTimeRecorder = guest.VisitTimeRecorder[index:]
	return ret

}

func (guest Guest) PackageCheck() {

	return

}

func CheckIP(addr string) {

	addr = net.ParseIP(addr).String()

	Mutex.Lock()
	defer Mutex.Unlock()

	guest, ok := GuestMap[addr]
	if !ok {
		guest = new(Guest)
		guest.GuestIP = addr
		GuestMap[addr] = guest
	}

	guest.visited()

	if guest.count() >= 5 {
		fillter.AppendBlackList(addr)
	}

}

func PackageLengthRecorder(addr string, length int) {

	Mutex.Lock()
	defer Mutex.Unlock()

	addr = net.ParseIP(addr).String()

	guest, ok := GuestMap[addr]
	if !ok {
		return
	}

	guest.sent(length)

	if tem := len(guest.PackageLengthRecoder); tem > 100 {
		guest.PackageLengthRecoder = guest.PackageLengthRecoder[tem-100:]
	}

}
