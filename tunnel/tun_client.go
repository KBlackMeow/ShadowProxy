package tunnel

import (
	"encoding/json"
	"net"
	"shadowproxy/cryptotools"
	"shadowproxy/logger"
)

type TunnelClient struct {
	ServiceAddr     string
	ServiceListener *net.UDPConn
	TunnelAddr      string
	Tunnels         map[uint32]*Tunnel
}

func (client TunnelClient) Run() {
	udpAddr, _ := net.ResolveUDPAddr("udp4", client.ServiceAddr)
	listener, err := net.ListenUDP("udp4", udpAddr)

	logger.Log("TUN", client.ServiceAddr, "tunnel connected.")

	if err != nil {
		logger.Error("TUN", err)
		return
	}
	client.ServiceListener = listener
	defer client.ServiceListener.Close()

	client.CreateTCPTunnel()

}

func (client TunnelClient) CreateTCPTunnel() {
	addr := "127.0.0.1:10001"
	udpAddr, _ := net.ResolveUDPAddr("udp4", client.TunnelAddr)
	tun := Tunnel{
		TunnelID:   uint32(cryptotools.EasyHash_uint64(addr)),
		TunnelConn: client.ServiceListener,
		TunnelAddr: udpAddr,
		TargetAddr: addr,
		Lines:      map[uint32]*Line{},
	}

	client.Tunnels[tun.TunnelID] = &tun
	tun.NewTun()

	for {
		buffer := make([]byte, 4096)
		n1, err := tun.TunnelConn.Read(buffer)
		if err != nil {
			logger.Error("TUN", err)
			continue
		}
		logger.Log("TUN", "client read ", n1)
		pkg := TunnelPackage{}
		e1 := json.Unmarshal(buffer[:n1], &pkg)
		if e1 != nil {
			logger.Error("TUN", e1)
			continue
		}

		if pkg.NewLine == 1 {

			conn, err := net.Dial("tcp", addr)
			if err != nil {
				logger.Error("TUN", err)
				return
			}

			line := Line{
				LineID: pkg.LineID,
				Tun:    tun,
				Conn:   conn,
			}

			line.NewLine()
			continue
		}

		if pkg.NewLine == 2 {
			continue
		}

		tun, ok := client.Tunnels[pkg.TunnelID]

		if !ok {
			continue
		}
		tun.SendToReal(pkg)

	}
}

func TunnelInit1() {
	tunnelClient := TunnelClient{
		ServiceAddr: "127.0.0.1:65533",
		TunnelAddr:  "127.0.0.1:65534",
		Tunnels:     map[uint32]*Tunnel{},
	}

	go tunnelClient.Run()
}
