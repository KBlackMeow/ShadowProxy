package tunnel

import (
	"encoding/json"
	"net"
)

type TunnelPackage struct {
	TunnelID  uint32
	LineID    uint32
	NewTun    uint32
	NewLine   uint32
	CloseTun  uint32
	CloseLine uint32
	Length    uint32
	Bytes     []byte
}

type Tunnel struct {
	TunnelID   uint32
	RemoteAddr *net.UDPAddr
	ListenConn net.Conn
	TunnelConn *net.UDPConn
}

func (tun Tunnel) Write(data []byte) (int, error) {
	n1, err := tun.TunnelConn.WriteToUDP(data, tun.RemoteAddr)
	return n1, err
}

func (tun Tunnel) Send(data []byte) (int, error) {
	n1, err := tun.ListenConn.Write(data)
	return n1, err
}

func (tun Tunnel) CloseTun() {
	pkg := TunnelPackage{
		TunnelID:  tun.TunnelID,
		CloseLine: 0,
		NewTun:    0,
		Length:    0,
		NewLine:   0,
		CloseTun:  1,
		Bytes:     []byte{},
	}
	data, _ := json.Marshal(pkg)
	tun.Write(data)
}

func (tun Tunnel) NewTun() {
	pkg := TunnelPackage{
		TunnelID:  tun.TunnelID,
		CloseLine: 0,
		NewTun:    1,
		Length:    0,
		NewLine:   0,
		CloseTun:  0,
		Bytes:     []byte{},
	}
	data, _ := json.Marshal(pkg)
	tun.Write(data)
}
