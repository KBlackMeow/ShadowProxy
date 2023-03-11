package proxy

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"net"
	"shadowproxy/config"
	"shadowproxy/cryptotools"
	"shadowproxy/logger"
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

		if config.ShadowProxyConfig.ReverseCrypt {
			go connections(backConn, linkConn, 0)
			go connections(linkConn, backConn, 1)
		} else {
			go connection(backConn, linkConn)
			go connection(linkConn, backConn)
		}

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

	if config.ShadowProxyConfig.ReverseCrypt {
		go connections(conn, linkConn, 0)
		go connections(linkConn, conn, 1)
	} else {
		go connection(conn, linkConn)
		go connection(linkConn, conn)
	}

}

func connections(from net.Conn, to net.Conn, tag int) {
	defer from.Close()
	defer to.Close()
	if tag == 1 {
		for {

			buffer := make([]byte, 4)

			var pkn uint32
			_, err := from.Read(buffer)
			if err != nil {
				return
			}

			err = binary.Read(bytes.NewReader(buffer[:4]), binary.BigEndian, &pkn)
			if err != nil {
				return
			}

			buffer = make([]byte, pkn)
			n, err := from.Read(buffer)
			if err != nil {
				return
			}

			var buff bytes.Buffer
			buff.Write(buffer[:n])
			for uint32(n) < pkn {
				tbuf := make([]byte, pkn-uint32(n))
				tn, err := from.Read(tbuf)
				if err != nil {
					return
				}
				buff.Write(tbuf)
				n += tn
			}

			buffer = cryptotools.Ase256Decode(buff.Bytes(), "12345678901234567890123456789012", "1234567890123456")
			_, err = to.Write(buffer)
			if err != nil {
				return
			}
		}
	} else if tag == 0 {
		for {
			buffer := make([]byte, 4096)
			n1, err := from.Read(buffer)
			if err != nil {
				return
			}

			var lengthBuf bytes.Buffer
			// err = binary.Write(&lengthBuf, binary.BigEndian, uint32(n1))
			// if err != nil {
			// 	return
			// }
			// to.Write(lengthBuf.Bytes())

			buffer = cryptotools.Ase256Encode(buffer[:n1], "12345678901234567890123456789012", "1234567890123456", aes.BlockSize)

			lengthBuf.Reset()
			err = binary.Write(&lengthBuf, binary.BigEndian, uint32(len(buffer)))
			if err != nil {
				return
			}
			to.Write(lengthBuf.Bytes())

			_, err = to.Write(buffer)
			if err != nil {
				return
			}
		}
	}

}

func connection(from net.Conn, to net.Conn) {
	defer from.Close()
	defer to.Close()
	for {
		buffer := make([]byte, 4096*16)
		n1, err := from.Read(buffer)
		if err != nil {
			logger.Error("REV", err)
			from.Close()
			to.Close()
			return
		}

		logger.Log("REV", from.RemoteAddr().String(), "->", to.RemoteAddr().String(), n1, "Bytes")

		_, err = to.Write(buffer[:n1])
		if err != nil {
			logger.Error("REV", err)
			from.Close()
			to.Close()
			return
		}
	}
}
