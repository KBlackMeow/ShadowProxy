package proxy

import (
	"crypto/aes"
	"net"
	"shadowproxy/config"
	"shadowproxy/cryptotools"
	"shadowproxy/logger"
	"strings"
	"time"
)

type RevProxyServer struct {
	ServerAddr string
	LinkAddr   string
	LinkConn   chan net.Conn
}

func (server RevProxyServer) Run() {
	server.LinkConn = make(chan net.Conn, 2)
	go server.LinkController()
	listener, err := net.Listen("tcp", server.ServerAddr)
	if err != nil {
		logger.Error("REV SER", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("REV SER", err)
			return
		}
		go server.Controller(conn)
	}
}
func (server RevProxyServer) LinkController() {
	listener, err := net.Listen("tcp", server.LinkAddr)
	if err != nil {
		logger.Error("REV SER LINK", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		logger.Log("REV SER Recive Link", conn.RemoteAddr().String())
		if err != nil {
			logger.Error("REV SER LINK", err)
			return
		}
		server.LinkConn <- conn
	}

}
func (server RevProxyServer) Controller(conn net.Conn) {
	for {
		buff := make([]byte, 32)
		n, err := conn.Read(buff)
		// TEST
		buff = cryptotools.Ase256Decode(buff[:n], "12345678901234567890123456789012", "1234567890123456")
		if err != nil {
			logger.Error("REV SER CON ", err)
			continue
		}
		if buff[0] == byte(255) {

			addr, err := server.CreateBackendListener(conn, string(buff[1:n]))
			if err != nil {
				logger.Error("REV SER CON", err)
				continue
			}
			// TEST
			buff = cryptotools.Ase256Encode([]byte(addr), "12345678901234567890123456789012", "1234567890123456", aes.BlockSize)
			_, err = conn.Write(buff)
			// _, err = conn.Write([]byte(addr))
			if err != nil {
				logger.Error("REV SER CON", err)
				continue
			}
		}
	}
}

func (server RevProxyServer) CreateBackendListener(conn net.Conn, backend string) (string, error) {

	logger.Log("REV SER BACK listen", backend)
	listener, err := net.Listen("tcp", backend)
	if err != nil {
		return "", err
	}
	go server.BackendListen(listener, conn)
	return backend, nil

}

func (server RevProxyServer) BackendListen(backend net.Listener, conn net.Conn) {
	defer backend.Close()
	for {
		backConn, err := backend.Accept()
		if err != nil {
			logger.Error("REV SER BACK", err)
			continue
		}
		buff := make([]byte, 16)
		buff[0] = 127
		// TEST
		buff = cryptotools.Ase256Encode(buff, "12345678901234567890123456789012", "1234567890123456", aes.BlockSize)

		_, err = conn.Write(buff)
		if err != nil {
			logger.Error("REV SER BACK", err)
			continue
		}

		linkConn := <-server.LinkConn
		go connection(backConn, linkConn, 0)
		go connection(linkConn, backConn, 1)
	}
}

type RevProxyClient struct {
	ServerAddr string
	LinkAddr   string
}

func (client RevProxyClient) Link(LocalAddr string, RemoteAddr string) {
	conn, err := net.Dial("tcp", client.ServerAddr)
	if err != nil {
		logger.Error("REV CLI", err)
		return
	}
	buff := make([]byte, 32)

	buff[0] = byte(255)
	copy(buff[1:], []byte(RemoteAddr))
	// TEST
	buff = cryptotools.Ase256Encode(buff, "12345678901234567890123456789012", "1234567890123456", aes.BlockSize)
	_, err = conn.Write(buff)
	if err != nil {
		logger.Error("REV CLI", err)
		return
	}
	go client.Controller(conn, LocalAddr)
}

func (client RevProxyClient) Controller(conn net.Conn, LocalAddr string) {
	for {
		buff := make([]byte, 4096)

		n, err := conn.Read(buff)

		// TEST
		buff = cryptotools.Ase256Decode(buff[:n], "12345678901234567890123456789012", "1234567890123456")
		if err != nil {
			logger.Error("REV CLI CON", err)
			return
		}
		if buff[0] == byte(127) {
			go client.Work(LocalAddr)
		} else {
			logger.Log("REV CLI INFO", string(buff[:n]))
		}

	}
}

func (client RevProxyClient) Work(LocalAddr string) {

	conn, err := net.Dial("tcp", LocalAddr)
	if err != nil {
		return
	}

	linkConn, err := net.Dial("tcp", client.LinkAddr)
	if err != nil {
		return
	}
	go connection(conn, linkConn, 0)
	go connection(linkConn, conn, 1)
}

// func connection(from net.Conn, to net.Conn, crypt int) {
// 	defer from.Close()
// 	defer to.Close()
// 	if crypt == 1 {
// 		for {
// 			buffer := make([]byte, 4096)
// 			n1, err := from.Read(buffer)
// 			if err != nil {
// 				return
// 			}

// 			fmt.Println(n1)
// 			buffer = cryptotools.Ase256Decode(buffer[:n1], "12345678901234567890123456789012", "1234567890123456")

// 			n := uint32(btoi(buffer[0:4]))
// 			fmt.Println("SER RECV ", n)
// 			buffer = buffer[4 : n+4]

// 			fmt.Println(1, "->", len(buffer), n)
// 			_, err = to.Write(buffer)
// 			if err != nil {
// 				return
// 			}
// 		}
// 	} else if crypt == 0 {
// 		for {
// 			buffer := make([]byte, 4080)
// 			n1, err := from.Read(buffer)
// 			if err != nil {
// 				return
// 			}

// 			// TEST
// 			send := make([]byte, 4096)
// 			copy(send[0:4], itob(uint32(n1)))
// 			copy(send[4:], buffer)
// 			buffer = send
// 			fmt.Println(len(buffer), n1)
// 			buffer = cryptotools.Ase256Encode(buffer[:n1+4], "12345678901234567890123456789012", "1234567890123456", aes.BlockSize)
// 			fmt.Println(0, "->", len(buffer), n1)

// 			_, err = to.Write(buffer)
// 			if err != nil {
// 				return
// 			}
// 		}
// 	}

// }

// func itob(i uint32) []byte {
// 	bt := make([]byte, 4)
// 	bt[0] = byte(i & 0xff)
// 	bt[1] = byte(i >> 8 & 0xff)
// 	bt[2] = byte(i >> 16 & 0xff)
// 	bt[3] = byte(i >> 24 & 0xff)
// 	fmt.Println("itob", bt)
// 	return bt
// }

// func btoi(bt []byte) uint32 {
// 	ret := uint32(0)
// 	fmt.Println("btoi", bt)
// 	ret = ret + uint32(bt[0]) + uint32(bt[1])<<8 + uint32(bt[2])<<16 + uint32(bt[3])<<24
// 	return ret
// }

func connection(from net.Conn, to net.Conn, crypt int) {
	defer from.Close()
	defer to.Close()
	for {
		buffer := make([]byte, 4096)
		n1, err := from.Read(buffer)
		if err != nil {
			return
		}
		_, err = to.Write(buffer[:n1])
		if err != nil {
			return
		}
	}
}

func RunRev() {
	server := RevProxyServer{
		ServerAddr: config.ShadowProxyConfig.ReverseServer,
		LinkAddr:   config.ShadowProxyConfig.ReverseLinkServer,
	}
	go server.Run()

	client := RevProxyClient{
		ServerAddr: config.ShadowProxyConfig.ReverseServer,
		LinkAddr:   config.ShadowProxyConfig.ReverseLinkServer,
	}
	time.Sleep(time.Second * 1)
	for _, v := range config.ShadowProxyConfig.ReverseRule {
		addrs := strings.Split(v, "->")

		go client.Link(addrs[0], addrs[1])

	}
}
