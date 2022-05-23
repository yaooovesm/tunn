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
	application.New()
	//app.Run()
	admin := administration.NewClientAdmin(config.Current.Admin)
	admin.Run()
}
