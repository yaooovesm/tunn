package tunnel

import (
	log "github.com/cihub/seelog"
	"tunn/config"
	"tunn/device"
	"tunn/traffic"
)

//
// AuthClientHandler
// @Description:
//
type AuthClientHandler struct {
	Client *Client
}

//
// GetDevice
// @Description:
// @receiver h
// @return *water.Interface
//
func (h *AuthClientHandler) GetDevice() device.Device {
	return h.Client.IFace
}

//
// OnMessage
// @Description:
// @receiver h
// @param msg
//
func (h *AuthClientHandler) OnMessage(msg string) {
	log.Info("receive message from server : ", msg)
}

//
// OnReport
// @Description:
// @receiver h
// @param payload
//
func (h *AuthClientHandler) OnReport(payload []byte) {
	log.Info("receive report data from server : ", len(payload), " bytes")
}

//
// OnLogin
// @Description:
// @receiver h
// @param err
// @param key
//
func (h *AuthClientHandler) OnLogin(err error, key []byte) {
	if err == nil {
		//crypt
		h.Client.PK = key
		rxDecryptFP := traffic.GetDecryptFP(config.Current.DataProcess, key)
		if rxDecryptFP != nil {
			h.Client.RxFP.Register(rxDecryptFP, "rx_decrypt")
		}
		//get cipher processor
		txEncryptFP := traffic.GetEncryptFP(config.Current.DataProcess, config.Current.DataProcess.Key)
		if txEncryptFP != nil {
			h.Client.TxFP.Register(txEncryptFP, "tx_encrypt")
		}
	}
}

//
// OnLogout
// @Description:
// @receiver h
// @param err
//
func (h *AuthClientHandler) OnLogout(err error) {
	if err == nil {
		err = ErrLogout
	}
	h.setClientDown(err)
}

//
// OnDisconnect
// @Description:
// @receiver h
//
func (h *AuthClientHandler) OnDisconnect() {
	log.Info("disconnected...")
	h.setClientDown(ErrDisconnect)
}

//
// OnKick
// @Description:
// @receiver h
//
func (h *AuthClientHandler) OnKick() {
	h.setClientDown(ErrStoppedByServer)
}

//
// OnRestart
// @Description:
// @receiver h
//
func (h *AuthClientHandler) OnRestart() {
	h.setClientDown(ErrRestartByServer)
}

//
// setClientDown
// @Description:
// @receiver h
// @param err
//
func (h *AuthClientHandler) setClientDown(err error) {
	h.Client.SetErr(err)
	h.Client.Stop()
	h.Client.multiConn.Close()
}
