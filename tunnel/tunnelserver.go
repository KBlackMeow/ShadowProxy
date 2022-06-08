package tunnel

import (
	"net"
	"shadowproxy/logger"
)

type TunServer struct {
	listener string
}

func (server TunServer) Run() {
	addr, _ := net.ResolveUDPAddr("udp4", server.listener)
	listener, err := net.ListenUDP("udp4", addr)

	if err != nil {
		logger.Error("Tunnel Server", err)
		return
	}

	defer listener.Close()
	logger.Log("Tunnel Server", server.listener, "tunnel server started.")

	for {

		buffer := make([]byte, 4096)
		n1, addr, err := listener.ReadFromUDP(buffer)

		if err != nil {
			logger.Error("Tunnel Server", err)
			return
		}

		pkg := GetTunPkgFromBytes(buffer, n1)

		// logger.Log("tunnel server", pkg)

		if pkg.flag == 0 {
			tun := Tunnel{server: server.listener, client: addr.String(), key: "123456"}
			logger.Log("Tunnel Server", tun)
		}

	}

}
