package tunnel

import (
	"errors"
	"strings"
	"tunn/authentication"
)

var (
	ErrDisconnect             = errors.New("disconnect")
	ErrDisconnectAccidentally = errors.New("disconnect accidentally")
	ErrLogin                  = errors.New("failed to login to server")
	ErrStoppedByServer        = errors.New("client stopped")
)

var (
	allowedErr = map[error]int{
		ErrDisconnect:                   1,
		ErrLogin:                        1,
		authentication.ErrAuthTimeout:   1,
		authentication.ErrAuthConnect:   1,
		authentication.ErrConnectFailed: 1,
		ErrDisconnectAccidentally:       1,
	}
	allowedErrStr = []string{
		"use of closed network connection",
	}
)

//
// IsAllowRestart
// @Description:
// @param err
// @param restart
// @return bool
//
func IsAllowRestart(err error, restart bool) bool {
	if !restart {
		return false
	}
	if _, ok := allowedErr[err]; ok {
		return true
	}
	s := err.Error()
	for i := range allowedErrStr {
		if strings.Contains(s, allowedErrStr[i]) {
			return true
		}
	}
	return false
}
