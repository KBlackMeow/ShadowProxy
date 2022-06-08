package tunnel

import (
	"bytes"
	"encoding/binary"
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
	key    string
}

func Run() {
	// x := TunPkg{1, 2, []byte("012")}
	// fmt.Println(x.toBytes())
	// byts, n := x.toBytes()
	// fmt.Println(GetTunPkgFromBytes(byts, n))

	ser := TunServer{listener: "0.0.0.0:11111"}
	go ser.Run()
	time.Sleep(time.Duration(100) * time.Millisecond)
	cli := TunClient{server: "127.0.0.1:11111"}

	cli.Link(2222, 3333)

}
