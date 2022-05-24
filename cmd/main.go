package main

import (
	"tunn/administration"
	"tunn/application"
	"tunn/config"
	"tunn/logging"
)

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
		app.Run()
	}
	admin := administration.NewClientAdmin(config.Current.Admin)
	admin.Run()
}
