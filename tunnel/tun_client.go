package tunnel

import (
	"encoding/json"
	"net"
	"shadowproxy/cryptotools"
	"shadowproxy/logger"
)

type TunnelClient struct {
	ServiceAddr     string
	ServiceListener net.Conn
	Tunnels         map[uint32]*Tunnel
}

func (client TunnelClient) Run() {
	udpAddr, _ := net.ResolveUDPAddr("udp4", client.ServiceAddr)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		logger.Error(err)
		return
	}
	client.ServiceListener = conn
	defer client.ServiceListener.Close()
	addr := "127.0.0.1:10001"
	tun := Tunnel{
		TunnelID:   uint32(cryptotools.EasyHash_uint64(addr)),
		TunnelConn: conn,
		TargetAddr: addr,
	}

	client.Tunnels[tun.TunnelID] = &tun
	tun.NewTun()

	for {
		buffer := make([]byte, 4096)
		n1, err := tun.TunnelConn.Read(buffer)
		if err != nil {
			logger.Error(err)
			continue
		}

		pkg := TunnelPackage{}
		e1 := json.Unmarshal(buffer[:n1], &pkg)
		if e1 != nil {
			logger.Error(e1)
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
		tun.Send(pkg)

	}

}
