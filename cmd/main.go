package main

import (
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
	app.Run()
}
