package tunnel

import (
	"encoding/json"
	"net"
	"shadowproxy/cryptotools"
	"shadowproxy/filter"
	"shadowproxy/ids"
	"shadowproxy/logger"
)

type TunnelServer struct {
	ServiceAddr     string
	ServiceListener *net.UDPConn
	Tunnels         map[uint32]*Tunnel
}

func (server TunnelServer) Run() {
	udpAddr, _ := net.ResolveUDPAddr("udp4", server.ServiceAddr)
	listener, err := net.ListenUDP("udp4", udpAddr)

	if err != nil {
		logger.Error("TUN", err)
		return
	}

	defer listener.Close()
	logger.Log("TUN", server.ServiceAddr, "tunnel started.")
	server.ServiceListener = listener

	for {

		buffer := make([]byte, 4096)
		n1, addr, err := server.ServiceListener.ReadFromUDP(buffer)

		if err != nil {
			logger.Error("TUN", err)
			return
		}

		if filter.Filter(addr.String()) {
			logger.Warn("TUN", addr.String(), "Alice is filtrated")
			continue
		}

		tunpak := TunnelPackage{}

		e1 := json.Unmarshal(buffer[0:n1], &tunpak)
		if e1 != nil {
			logger.Error(e1)
			continue
		}
		go ids.PackageLengthRecorder(addr.String(), n1)

		if tunpak.NewTun == 1 {
			server.CreateTCPTunnel(addr)
			continue
		}
		if tunpak.NewTun == 2 {
			server.CreateUDPTunnel(addr)
			continue
		}

		tun, ok := server.Tunnels[tunpak.TunnelID]
		if !ok {
			continue
		}
		tun.Send(tunpak)
	}

}

func (server TunnelServer) CreateTCPTunnel(remoteAddr *net.UDPAddr) {

	addr := "0.0.0.0:44556"
	listener, err := net.Listen("tcp4", addr)

	if err != nil {
		logger.Error("TUN", err)
		return
	}
	tun := Tunnel{
		ListenConn: server.ServiceListener,
		TunnelID:   uint32(cryptotools.EasyHash_uint64(remoteAddr.String())),
		RemoteAddr: remoteAddr,
		TunnelConn: server.ServiceListener,
	}

	server.Tunnels[tun.TunnelID] = &tun
	defer listener.Close()
	logger.Log("TUN", addr, "tunnel started.")
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("TUN", err)
			return
		}

		line := Line{
			Tun:    tun,
			Conn:   conn,
			LineID: uint32(cryptotools.EasyHash_uint64(conn.RemoteAddr().String())),
		}
		line.NewLine()

	}
}

func (server TunnelServer) CreateUDPTunnel(remoteAddr *net.UDPAddr) {

}
