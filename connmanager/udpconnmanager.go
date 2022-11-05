package connmanager

import (
	"net"
	"shadowproxy/logger"
	"sync"
	"time"
)

type UDPConn struct {
	Addr         *net.UDPAddr
	BackendConn  net.Conn
	ListenerConn *net.UDPConn
	TTL          int64
	RecvTime     time.Time
}

func (conn UDPConn) Write(buff []byte) (int, error) {

	return conn.ListenerConn.WriteToUDP(buff, conn.Addr)

}

func (conn UDPConn) Close() {

	conn.BackendConn.Close()

}

func (conn *UDPConn) Backword() {
	from := conn.BackendConn
	to := conn
	for {
		buffer := make([]byte, 4096)
		n1, err := from.Read(buffer)

		if err != nil {
			logger.Error("UDP", err)
			conn.Close()
			return
		}

		n2, err := to.Write(buffer[:n1])
		if err != nil {
			logger.Error("UDP", err)
			conn.Close()
			return
		}

		logger.Log("UDP", from.RemoteAddr().String(), "->", to.Addr.String(), n2, "Bytes")
		to.RecvTime = time.Now()
	}
}

var UDPConns = map[string]*UDPConn{}
var UDPMutex = new(sync.Mutex)

func SetUDPConn(addr string, conn *UDPConn) {

	UDPMutex.Lock()
	defer UDPMutex.Unlock()
	UDPConns[addr] = conn

}

func GetUDPConn(addr string) (*UDPConn, bool) {

	UDPMutex.Lock()
	defer UDPMutex.Unlock()
	udpConn, ok := UDPConns[addr]
	return udpConn, ok

}

func CloseUDPConnFromAddr(addr string) {
	UDPMutex.Lock()
	defer UDPMutex.Unlock()
	UDPConns[addr].BackendConn.Close()
	delete(UDPConns, addr)
}

func AutoUDPConnClose() {

	for {
		for k, v := range UDPConns {
			if time.Since(v.RecvTime).Milliseconds() > v.TTL {
				v.Close()
				delete(UDPConns, k)
			}
		}

		time.Sleep(time.Duration(500) * time.Millisecond)
	}

}

func AllUDPConnClose() {

	for k, v := range UDPConns {
		v.Close()
		delete(UDPConns, k)
	}

}

func init() {

	go AutoUDPConnClose()

}
