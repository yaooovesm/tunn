package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	Current.ReadFromFile("D:\\code\\golang\\tunnel\\common\\config\\tunnel_client.json")
	marshal, err := json.Marshal(Current)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}

func TestConfigRead(t *testing.T) {
	cfg := Config{
		Global: Global{
			Tunnel: Tunnel{
				Address:  "",
				Port:     0,
				Protocol: "",
			},
			MTU:       0,
			Pprof:     0,
			MultiConn: 0,
		},
		User: User{
			Account:  "",
			Password: "",
		},
		Routes: []Route{},
		Device: Device{
			CIDR: "",
			DNS:  "",
		},
		Auth: Auth{
			Address: "",
			Port:    0,
		},
		DataProcess: DataProcess{
			CipherType: "",
		},
		Security: Security{
			CertPem: "",
		},
	}
	marshal, err := json.Marshal(cfg)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}
