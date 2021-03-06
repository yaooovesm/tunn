package tunnel

import (
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"net"
	"net/url"
	"time"
	"tunn/config"
	"tunn/transmitter"
)

//
// WSClientHandler
// @Description:
//
type WSClientHandler struct {
	url url.URL
}

//
// AfterInitialize
// @Description:
// @receiver h
// @param client
//
func (h *WSClientHandler) AfterInitialize(client *Client) {
	u := url.URL{Scheme: "ws", Host: client.Address, Path: "/" + client.AuthClient.WSKey + "/access_point"}
	h.url = u
	log.Info("connect to ws server : ", h.url.String())
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
func (h *WSClientHandler) CreateAndSetup(address string, config config.Config) (conn net.Conn, err error) {
	dialer := websocket.Dialer{
		HandshakeTimeout:  time.Second * time.Duration(45),
		EnableCompression: false,
	}
	wsconn, _, err := dialer.Dial(h.url.String(), nil)

	return transmitter.WrapWSConn(wsconn), err
}
