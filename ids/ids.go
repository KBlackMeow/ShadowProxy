package ids

import (
	"shadowproxy/fillter"
	"strings"
	"time"
)

var GuestMap = map[string]*Guest{}

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

func CheckAddr(addr string) {
	if tem := strings.Split(addr, ":"); len(tem) == 2 {
		addr = tem[0]
	}

	guest, ok := GuestMap[addr]
	if !ok {
		guest = new(Guest)
		guest.GuestIP = addr
		GuestMap[addr] = guest
	}

	guest.visited()

	if guest.count() > 5 {
		fillter.AppendBlackList(addr)
	}
}

func PackageLengthRecorder(addr string, length int) {

	if tem := strings.Split(addr, ":"); len(tem) == 2 {
		addr = tem[0]
	}

	guest, ok := GuestMap[addr]
	if !ok {
		return
	}

	guest.sent(length)

	// logger.Log(addr, length, len(guest.PackageLengthRecoder))

	if tem := len(guest.PackageLengthRecoder); tem > 100 {
		guest.PackageLengthRecoder = guest.PackageLengthRecoder[tem-100:]
	}

}
