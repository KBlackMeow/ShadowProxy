package tunnel

type TunClient struct {
	RemoteAddr string
}

func (client TunClient) Link(src uint16, dst uint16, flag uint16) {

	tun := Tunnel{server: client.RemoteAddr, client: "127.0.0.1", src: src, dst: dst, flag: flag, key: "123456"}
	tun.connect()
	// server, err := net.Dial("udp", client.RemoteAddr)
	// if err != nil {
	// 	logger.Error("UDP", err)
	// 	return
	// }
	// defer server.Close()

	// pkg := TunPkg{src: uint16(src), dst: uint16(dst), flag: 0, pkg: []byte{}}
	// byts, _ := pkg.toBytes()

	// _, e1 := server.Write(byts)
	// if e1 != nil {
	// 	logger.Error("Tunnel Client", e1)
	// 	return
	// }

	// buffer := make([]byte, 4096)
	// n, e2 := server.Read(buffer)
	// if e2 != nil {
	// 	logger.Error("Tunnel Client", e2)
	// 	return
	// }
	// if GetTunPkgFromBytes(buffer, n).flag == 0 {
	// 	tun := Tunnel{server: server.RemoteAddr().String(), client: server.LocalAddr().String(), key: "123456"}
	// 	logger.Log("Tunnel Client", tun)
	// }

}
