package tunnel

import (
	"encoding/json"
	"net"
	"shadowproxy/cryptotools"
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

		logger.Log("TUN", addr.String(), "->", server.ServiceAddr, n1, "Bytes")

		tunpak := TunnelPackage{}

		e1 := json.Unmarshal(buffer[0:n1], &tunpak)

		if e1 != nil {
			logger.Error("TUN", e1)
			continue
		}

		if tunpak.NewTun == 1 {
			go server.CreateTCPTunnel(addr, tunpak.TunnelID)
			continue
		}
		if tunpak.NewTun == 2 {
			go server.CreateUDPTunnel(addr, tunpak.TunnelID)
			continue
		}

		tun, ok := server.Tunnels[tunpak.TunnelID]
		if !ok {
			continue
		}
		go tun.SendToReal(tunpak)
	}

}

func (server TunnelServer) CreateTCPTunnel(remoteAddr *net.UDPAddr, tunid uint32) {

	addr := "0.0.0.0:10002"
	listener, err := net.Listen("tcp4", addr)

	if err != nil {
		logger.Error("TUN", err)
		return
	}

	logger.Log("TUN", addr, "target server listening.")

	tun := Tunnel{
		TunnelID:   tunid,
		TunnelAddr: remoteAddr,
		TunnelConn: server.ServiceListener,
		Lines:      map[uint32]*Line{},
	}

	server.Tunnels[tun.TunnelID] = &tun
	defer listener.Close()
	logger.Log("TUN", addr, "tunnel", tun.TunnelID, " started.")
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

func (server TunnelServer) CreateUDPTunnel(remoteAddr *net.UDPAddr, tunid uint32) {

}

func TunnelInit2() {
	tunserver := TunnelServer{
		ServiceAddr: "0.0.0.0:65534",
		Tunnels:     map[uint32]*Tunnel{},
	}

	go tunserver.Run()
}
