package tunnel

import (
	log "github.com/cihub/seelog"
	"github.com/xtaci/kcp-go"
	"net"
	"time"
	"tunn/config"
)

//
// KCPClientHandler
// @Description:
//
type KCPClientHandler struct {
}

//
// AfterInitialize
// @Description:
// @receiver h
// @param client
//
func (h KCPClientHandler) AfterInitialize(client *Client) {

}

//
// CreateAndSetup
// @Description:
// @receiver h
// @param address
// @return conn
// @return err
//
func (h *KCPClientHandler) CreateAndSetup(address string, config config.Config) (conn net.Conn, err error) {
	session, err := kcp.DialWithOptions(address, nil, 10, 3)
	if err != nil {
		return nil, err
	}
	//保持连接
	err = session.SetDeadline(time.Time{})
	if err != nil {
		log.Info("set kcp deadline failed : ", err)
	}
	return session, nil
}
