package networking

import (
	"fmt"
	"testing"
	"time"
	"tunn/config"
	"tunn/device"
)

func TestSystemrt(t *testing.T) {
	dev, err := device.NewTunDeviceWithConfig(config.Config{
		Global: config.Global{MTU: 1500},
		Device: config.Device{
			CIDR: "10.0.0.77/24",
			DNS:  "223.5.5.5",
		},
	})
	if err != nil {
		fmt.Println("create : ", err)
		return
	}
	err = dev.Setup()
	if err != nil {
		fmt.Println("setup : ", err)
		return
	}
	AddSystemRoute("172.24.8.0/24", dev.Name())
	time.Sleep(time.Second * 5)
	err = dev.Close()
	if err != nil {
		fmt.Println("close : ", err)
		return
	}
	time.Sleep(time.Second * 10)
}
