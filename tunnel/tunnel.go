package tunnel

import (
	"net"
	"tunn/config"
)

//
// ClientConnHandler
// @Description:
//
type ClientConnHandler interface {
	//
	// AfterInitialize
	// @Description:
	// @param client
	//
	AfterInitialize(client *Client)
	//
	// CreateAndSetup
	// @Description:
	// @param Address
	// @return conn
	// @return err
	//
	CreateAndSetup(address string, config config.Config) (conn net.Conn, err error)
}
