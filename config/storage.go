package config

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"os"
)

//
// ClientConfigStorage
// @Description:
//
type ClientConfigStorage struct {
	User     User     `json:"user"`
	Auth     Auth     `json:"auth"`
	Security Security `json:"security"`
	Admin    Admin    `json:"admin"`
}

//
// ReadFromFile
// @Description:
// @receiver cfg
// @param path
//
func (cfg *ClientConfigStorage) ReadFromFile(path string) {
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
	_ = json.Unmarshal(bytes, cfg)
}

//
// ToConfig
// @Description:
// @receiver cfg
// @return Config
//
func (cfg *ClientConfigStorage) ToConfig() Config {
	return Config{
		User:     cfg.User,
		Auth:     cfg.Auth,
		Security: cfg.Security,
		Admin:    cfg.Admin,
	}
}

//
// Dump
// @Description:
// @receiver cfg
//
func (cfg *ClientConfigStorage) Dump(path string) error {
	//检查是否保存了密码
	currentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	current := ClientConfigStorage{}
	err = json.Unmarshal(currentBytes, &current)
	if err != nil {
		return err
	}
	sto := cfg
	if current.User.Password != "" {
		sto.User.Password = current.User.Password
	}
	bytes, err := json.MarshalIndent(sto, "", "    ")
	if err != nil {
		return err
	}
	log.Info("config dump to : ", path)
	return ioutil.WriteFile(path, bytes, 0600)
}
