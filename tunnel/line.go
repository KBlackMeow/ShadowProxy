package tunnel

import (
	"encoding/json"
	"net"
	"shadowproxy/logger"
)

type Line struct {
	LineID uint32
	Tun    Tunnel
	Conn   net.Conn
}

func (line Line) ListenFromLine() {

	buffer := make([]byte, 4096)

	for {
		n1, err := line.Conn.Read(buffer)
		if err != nil {
			line.CloseLine()
			return
		}

		n2, err := line.Tun.Write(buffer[:n1])
		logger.Log("TUN", line.Conn.RemoteAddr().String(), "->", line.Tun.TunnelAddr.String(), n2, "Bytes")

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
	return line.Tun.Write(data)
}

func (line Line) SendToLine(byt []byte) (int, error) {
	return line.Conn.Write(byt)
}

func (line Line) NewLine() {

	_, ok := line.Tun.Lines[line.LineID]

	if ok {
		return
	}

	line.Tun.Lines[line.LineID] = &line

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

	go line.ListenFromLine()

	logger.Log("TUN", "Tunnel ", line.Tun.TunnelID, "Line", line.LineID, "connected.")
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
	delete(line.Tun.Lines, line.LineID)

}
