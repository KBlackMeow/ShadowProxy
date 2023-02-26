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
		tun.Send(tunpak.Bytes)
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

		go line.Listen()
	}
}

func (server TunnelServer) CreateUDPTunnel(remoteAddr *net.UDPAddr) {

}

type Line struct {
	LineID uint32
	Tun    Tunnel
	Conn   net.Conn
}

func (line Line) Listen() {

	buffer := make([]byte, 4096)

	for {
		n1, err := line.Tun.ListenConn.Read(buffer)
		if err != nil {
			line.CloseLine()
			return
		}

		n2, err := line.Tun.Write(buffer[:n1])
		logger.Log("TUN", line.Conn.RemoteAddr().String(), "->", line.Tun.RemoteAddr.String(), n2, "Bytes")

		if err != nil {
			line.CloseLine()
			return
		}
	}
}

func (line Line) WriteToLine(byt []byte) (int, error) {
	pkg := TunnelPackage{
		TunnelID:  line.Tun.TunnelID,
		LineID:    line.LineID,
		CloseLine: 0,
		CloseTun:  0,
		NewTun:    0,
		NewLine:   0,
		Length:    uint32(len(byt)),
		Bytes:     byt,
	}

	data, _ := json.Marshal(pkg)
	n1, err := line.Tun.Write(data)
	return n1, err
}

func (line Line) NewLine() {
	pkg := TunnelPackage{
		TunnelID:  line.Tun.TunnelID,
		LineID:    line.LineID,
		CloseLine: 0,
		NewTun:    0,
		NewLine:   1,
		Length:    0,
		Bytes:     []byte{},
	}
	data, _ := json.Marshal(pkg)
	line.Tun.Write(data)
}

func (line Line) CloseLine() {
	pkg := TunnelPackage{
		TunnelID:  line.Tun.TunnelID,
		LineID:    line.LineID,
		CloseLine: 1,
		NewTun:    0,
		Length:    0,
		NewLine:   0,
		Bytes:     []byte{},
	}
	data, _ := json.Marshal(pkg)
	line.Tun.Write(data)
}
