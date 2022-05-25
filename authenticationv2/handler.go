package authenticationv2

import (
	"tunn/device"
)

//
// AuthClientHandler
// @Description:
//
type AuthClientHandler interface {
	//
	// GetDevice
	// @Description:
	// @return *Device
	//
	GetDevice() device.Device
	//
	// OnMessage
	// @Description:
	// @param msg
	//
	OnMessage(msg string)
	//
	// OnReport
	// @Description:
	// @param payload
	//
	OnReport(payload []byte)
	//
	// OnLogin
	// @Description:
	// @param reply
	//
	OnLogin(err error, key []byte)
	//
	// OnLogout
	// @Description:
	// @param reply
	//
	OnLogout(err error)
	//
	// OnDisconnect
	// @Description:
	//
	OnDisconnect()
	//
	// OnKick
	// @Description:
	//
	OnKick()
	//
	// OnRestart
	// @Description:
	//
	OnRestart()
}
