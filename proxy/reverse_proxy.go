package proxy

import (
	"net"
	"shadowproxy/config"
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
	buff := make([]byte, 24)
	for {
		n, err := conn.Read(buff)
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
			_, err = conn.Write([]byte(addr))
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
		buff := make([]byte, 1)
		buff[0] = 127
		_, err = conn.Write(buff)
		if err != nil {
			logger.Error("REV SER BACK", err)
			continue
		}

		linkConn := <-server.LinkConn
		go connection(backConn, linkConn)
		go connection(linkConn, backConn)
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
	buff := make([]byte, 24)

	buff[0] = byte(255)
	copy(buff[1:], []byte(RemoteAddr))
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
	go connection(conn, linkConn)
	go connection(linkConn, conn)
}

func connection(from net.Conn, to net.Conn) {
	defer from.Close()
	defer to.Close()
	buffer := make([]byte, 4096)
	for {
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
