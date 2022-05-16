package networking

import (
	log "github.com/cihub/seelog"
	"os/exec"
)

var Initialized = false

/**
linux在暴露网络时需要开启内核转发，并且通过iptables进行地址伪装
临时开启内核转发
echo 1 > /proc/sys/net/ipv4/ip_forward
开启地址伪装
iptables -t nat -A  POSTROUTING -s [tunn内网] -j MASQUERADE
如 iptables -t nat -A  POSTROUTING -s 192.168.0.0/24 -j MASQUERADE
*/

//
// RouteSupport
// @Description:
//
func RouteSupport() {
	if Initialized {
		return
	}
	Initialized = true
}

//
// AddSystemRoute
// @Description:
// @param network
// @param dev
//
func AddSystemRoute(network string, dev string) error {
	log.Info("[", dev, "]add system route : ", network)
	err := command("/sbin/ip", "route", "add", network, "dev", dev)
	if err != nil {
		return log.Warn("import ", network, " failed : ", err)
	}
	return nil
}

//
// command
// @Description:
// @param c
// @param args
//
func command(c string, args ...string) error {
	cmd := exec.Command(c, args...)
	return cmd.Run()
}
