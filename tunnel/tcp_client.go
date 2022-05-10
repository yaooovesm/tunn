package tunnel

import (
	"net"
	"tunn/config"
)

//
// TCPClientHandler
// @Description:
//
type TCPClientHandler struct {
}

//
// AfterInitialize
// @Description:
// @receiver h
// @param client
//
func (h TCPClientHandler) AfterInitialize(client *Client) {

}

//
// CreateAndSetup
// @Description:
// @receiver h
// @param address
// @param config
// @return conn
// @return err
//
func (h *TCPClientHandler) CreateAndSetup(address string, config config.Config) (conn net.Conn, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	_ = tcpConn.SetKeepAlive(true)
	return tcpConn, nil
}
