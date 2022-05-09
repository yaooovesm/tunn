package authentication

import (
	"tunn/config"
	"tunn/transmitter"
)

//
// Connection
// @Description:
//
type Connection struct {
	UUID   string
	Config config.Config
	Tunn   *transmitter.Tunnel
}

//
// Disconnect
// @Description:
// @receiver c
//
func (c *Connection) Disconnect() {
	if c.Tunn != nil {
		_ = c.Tunn.Close()
	}
}
