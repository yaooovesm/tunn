package config

import (
	"encoding/json"
	"flag"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"os"
	"tunn/config/protocol"
)

var Location = ""

var Current = Config{}

//
// Config
// @Description:
//
type Config struct {
	Global      Global      `json:"global"`
	User        User        `json:"user"`
	Routes      []Route     `json:"route"`
	Device      Device      `json:"device"`
	Auth        Auth        `json:"auth"`
	DataProcess DataProcess `json:"data_process"`
	Security    Security    `json:"security"`
	Runtime     Runtime     `json:"runtime"`
	Admin       Admin       `json:"admin"`
	Limit       Limit       `json:"limit"`
}

//
// Global
// @Description: global config
//
type Global struct {
	Tunnel
	MTU       int `json:"mtu"`
	Pprof     int `json:"pprof"`
	MultiConn int `json:"multi_connection"`
}

//
// ReadFromFile
// @Description:
// @receiver cfg
// @param path
//
func (cfg *Config) ReadFromFile(path string) {
	if path == "" {
		_ = log.Error("config not specific")
		os.Exit(-1)
		return
	}
	log.Info("load config from : ", path)
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		_ = log.Error("failed to open config file : " + err.Error())
		os.Exit(-1)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		_ = log.Error("failed to read config file : " + err.Error())
		os.Exit(-1)
	}
	cfg.SetDefaultValue()
	_ = json.Unmarshal(bytes, cfg)
}

//
// SetDefaultValue
// @Description:
// @receiver cfg
//
func (cfg *Config) SetDefaultValue() {
	cfg.Global.Protocol = protocol.TCP
	cfg.Global.MTU = 1400
}

//
// Check
// @Description:
// @receiver cfg
//
func (cfg *Config) Check() {
	//if cfg.Global.Protocol == protocol.WSS || cfg.Global.Protocol == protocol.WS {
	//	log.Info("protocol ", cfg.Global.Protocol, " : multi_connection reset to 1")
	//	cfg.Global.MultiConn = 1
	//}
}

//
// MergePushed
// @Description:
// @receiver cfg
// @param push
//
func (cfg *Config) MergePushed(push PushedConfig) {
	cfg.Global.Address = push.Global.Address
	cfg.Global.Protocol = push.Global.Protocol
	cfg.Global.Port = push.Global.Port
	cfg.Global.MultiConn = push.Global.MultiConnection
	cfg.Global.MTU = push.Global.Mtu
	cfg.Routes = push.Routes
	cfg.Device = push.Device
	cfg.DataProcess.CipherType = push.DataProcess.CipherType
	cfg.Limit = push.Limit
}

//
// Clear
// @Description:
// @receiver cfg
//
func (cfg *Config) Clear() {
	cfg.Global.Address = ""
	cfg.Global.Protocol = ""
	cfg.Global.Port = 0
	cfg.Global.MultiConn = 0
	cfg.Global.MTU = 0
	cfg.Routes = []Route{}
	cfg.Device = Device{}
	cfg.DataProcess.CipherType = ""
}

//
// ToStorageModel
// @Description:
// @receiver cfg
// @return ClientConfigStorage
//
func (cfg *Config) ToStorageModel() ClientConfigStorage {
	storage := ClientConfigStorage{
		Auth:     cfg.Auth,
		Security: cfg.Security,
		Admin:    cfg.Admin,
		User: User{
			Account: cfg.User.Account,
		},
	}
	return storage
}

//
// Storage
// @Description:
// @receiver cfg
// @return error
//
func (cfg *Config) Storage(saveAccount bool) error {
	storage := cfg.ToStorageModel()
	if !saveAccount {
		storage.User = User{}
	}
	return storage.Dump(Location)
}

//
// Load
// @Description:
//
func Load() {
	c := flag.String("c", "", "config path")
	flag.Parse()
	Location = *c
	storage := ClientConfigStorage{}
	storage.ReadFromFile(Location)
	Current = storage.ToConfig()
	Current.Check()
	Current.Runtime.Collect()
}
