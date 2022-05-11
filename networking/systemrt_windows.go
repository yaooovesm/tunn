package networking

import (
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"net"
	"os/exec"
	"strconv"
)

var Initialized = false

/**
Windows在暴露网络时需要开启Routing and Remote Access服务
cmd | sc config RemoteAccess start=auto
cmd | sc start RemoteAccess
开启路由转发
powershell | reg add HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters /v IPEnableRouter /D 1 /f
*/

//
// RouteSupport
// @Description:
//
func RouteSupport() {
	if Initialized {
		return
	}
	log.Info("route support set on")
	err := command("cmd", "sc config RemoteAccess start=auto")
	if err != nil {
		_ = log.Warn("route support : ", err)
	}
	err = command("cmd", "sc start RemoteAccess")
	if err != nil {
		_ = log.Warn("route support : ", err)
	}
	err = command("PowerShell", "reg add HKLM\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters /v IPEnableRouter /D 1 /f")
	if err != nil {
		_ = log.Warn("route support : ", err)
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
	ip, ipNet, err := net.ParseCIDR(network)
	if err != nil {
		_ = log.Warn("import ", network, " failed : ", err)
		return err
	}
	//PowerShell route add -p [network] mask [mask] [dev_ip]
	devIp, index, err := getIpv4ByInterfaceName(dev)
	if err != nil {
		_ = log.Warn("import ", network, " failed : ", err)
		return err
	}
	err = command("PowerShell", "route", "add", ip.String(), "mask", ipv4MaskString(ipNet.Mask), devIp, "IF", strconv.Itoa(index))
	if err != nil {
		_ = log.Warn("import ", network, " failed : ", err)
		return err
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
	log.Debug("exec : ", cmd.String())
	return cmd.Run()
}

//
// ipv4MaskString
// @Description:
// @param m
// @return string
//
func ipv4MaskString(m []byte) string {
	if len(m) != 4 {
		panic("ipv4Mask: len must be 4 bytes")
	}
	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

//
// getIpv4ByInterfaceName
// @Description:
// @param name
// @return string
// @return error
//
func getIpv4ByInterfaceName(name string) (string, int, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", -1, err
	}
	for _, inter := range interfaces {
		if inter.Name == name {
			addrs, err := inter.Addrs()
			if err != nil {
				return "", -1, err
			}
			for i := range addrs {
				if isv4(addrs[i]) {
					ip, _, err := net.ParseCIDR(addrs[i].String())
					if err != nil {
						return "", -1, err
					}
					return ip.String(), inter.Index, nil
				}
			}
		}
	}
	return "", -1, errors.New("interface not found")
}

//
// isv4
// @Description:
// @param addr
// @return bool
//
func isv4(addr net.Addr) bool {
	ip := addr.String()
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return true
		case ':':
			return false
		}
	}
	return false
}
