package tunnel

import (
	"net"
	"shadowproxy/logger"
)

type TunClient struct {
	RemoteAddr string
}

func (client TunClient) Link(src uint16, dst uint16) {
	server, err := net.Dial("udp", client.RemoteAddr)
	if err != nil {
		logger.Error("UDP", err)
		return
	}
	defer server.Close()

	pkg := TunPkg{src: uint16(src), dst: uint16(dst), flag: 0, pkg: []byte{}}
	byts, _ := pkg.toBytes()

	_, e := server.Write(byts)

	if e != nil {
		logger.Error(e)
		return
	}

	tun := Tunnel{server: server.RemoteAddr().String(), client: server.LocalAddr().String(), key: "123456"}

	logger.Log("Tunnel Client", tun)

}
