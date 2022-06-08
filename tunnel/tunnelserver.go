package tunnel

import (
	"net"
	"shadowproxy/logger"
)

type TunServer struct {
	LocalAddr string
}

func (server TunServer) Run() {
	addr, _ := net.ResolveUDPAddr("udp4", server.LocalAddr)
	listener, err := net.ListenUDP("udp4", addr)

	if err != nil {
		logger.Error("Tunnel Server", err)
		return
	}

	defer listener.Close()
	logger.Log("Tunnel Server", server.LocalAddr, "tunnel server started.")

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
			byts, n := pkg.toBytes()
			listener.WriteToUDP(byts[:n], addr)

			tun := Tunnel{server: server.LocalAddr, client: addr.String(), key: "123456"}
			logger.Log("Tunnel Server", tun)
		}

	}

}
