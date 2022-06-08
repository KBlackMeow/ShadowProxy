package tunnel

import (
	"bytes"
	"encoding/binary"
	"net"
	"shadowproxy/logger"
	"time"
)

type TunPkg struct {
	src  uint16
	dst  uint16
	flag uint16
	pkg  []byte
}

func (tunPkg TunPkg) toBytes() ([]byte, int) {

	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, tunPkg.src)
	binary.Write(buff, binary.BigEndian, tunPkg.dst)
	binary.Write(buff, binary.BigEndian, tunPkg.flag)
	binary.Write(buff, binary.BigEndian, tunPkg.pkg)

	return buff.Bytes(), buff.Len()
}

func GetTunPkgFromBytes(msg []byte, n int) *TunPkg {

	pkg := new(TunPkg)
	buff := bytes.NewBuffer(msg[:n])
	binary.Read(buff, binary.BigEndian, &pkg.src)
	binary.Read(buff, binary.BigEndian, &pkg.dst)
	binary.Read(buff, binary.BigEndian, &pkg.flag)
	// binary.Read(buff, binary.BigEndian, pkg.pkg)
	pkg.pkg = buff.Bytes()
	return pkg

}

type Tunnel struct {
	server string
	client string
	src    uint16
	dst    uint16
	flag   uint16
	key    string
	conn   net.Conn
}

func (tun Tunnel) connect() {

	server, err := net.Dial("udp", tun.server)
	if err != nil {
		logger.Error("UDP", err)
		return
	}
	defer server.Close()

	pkg := TunPkg{src: tun.src, dst: tun.dst, flag: tun.flag, pkg: []byte{}}
	byts, _ := pkg.toBytes()

	_, e1 := server.Write(byts)
	if e1 != nil {
		logger.Error("Tunnel Client", e1)
		return
	}

	buffer := make([]byte, 4096)
	_, e2 := server.Read(buffer)
	if e2 != nil {
		logger.Error("Tunnel Client", e2)
		return
	}

	tun.conn = server

}

func (tun Tunnel) Write(buff []byte, n int) (int, error) {
	return tun.conn.Write(buff[:n])
}

func (tun Tunnel) Read(buff []byte) (int, error) {
	return tun.Read(buff)
}

func Run() {
	// x := TunPkg{1, 2, []byte("012")}
	// fmt.Println(x.toBytes())
	// byts, n := x.toBytes()
	// fmt.Println(GetTunPkgFromBytes(byts, n))

	ser := TunServer{LocalAddr: "0.0.0.0:11111"}
	go ser.Run()
	time.Sleep(time.Duration(500) * time.Millisecond)
	cli := TunClient{RemoteAddr: "127.0.0.1:11111"}

	cli.Link(2222, 3333, 1)

}
