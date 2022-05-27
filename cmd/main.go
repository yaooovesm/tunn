package main

import (
	"embed"
	"tunn/administration"
	"tunn/application"
	"tunn/config"
	"tunn/logging"
)

//go:embed static
var static embed.FS

//
// main
// @Description: entrance
//
func main() {
	//initialize log
	logging.Initialize()
	//load config
	config.Load()
	//create and run application
	app := application.New()
	if config.Current.User.Account != "" && config.Current.User.Password != "" {
		_ = app.Run()
	}
	admin := administration.NewClientAdmin(config.Current.Admin, static)
	admin.Run()
}
