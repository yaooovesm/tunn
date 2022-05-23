package application

import (
	log "github.com/cihub/seelog"
	"os"
	"time"
	"tunn/config"
	"tunn/config/protocol"
	"tunn/tunnel"
	"tunn/version"
)

var Current *Application

//
// Application
// @Description:
//
type Application struct {
	Config   config.Config
	Protocol protocol.Name
	serv     Service
	Running  bool
	Error    string
}

//
// New
// @Description:
// @return *Application
//
func New() *Application {
	app := &Application{
		Config:   config.Current,
		Protocol: config.Current.Global.Protocol,
		Running:  false,
	}
	Current = app
	return app
}

//
// Service
// @Description:
//
type Service interface {
	//
	// Init
	// @Description: run once before start service
	// @return error
	//
	Init() error
	//
	// Start
	// @Description: start service
	// @return error
	//
	Start() error
	//
	// Stop
	// @Description: stop service
	//
	Stop()
	//
	// Terminate
	// @Description:
	//
	Terminate()
}

//
// runService
// @Description:
// @receiver app
// @param serv
//
func (app *Application) runService(serv Service) {
	app.Error = ""
	log.Info("tunnel version : ", version.Version)
	if version.Develop {
		_ = log.Warn("当前版本为测试版本！")
	}
	//app.PProf()
	if serv == nil {
		_ = log.Warn("tunnel server type not support")
		os.Exit(-1)
		return
	}
	ch := make(chan error, 1)
	//TODO 解决重复初始化问题
	if err := serv.Init(); err != nil {
		_ = log.Warn("init failed : ", err)
		os.Exit(-1)
		return
	} else {
		log.Info("service init success...")
	}
	go func() {
		app.Running = true
		ch <- serv.Start()
	}()
	for {
		select {
		case err := <-ch:
			if err != nil {
				serv.Stop()
				_ = log.Warn("tunnel exit with error : ", err.Error())
				app.Error = err.Error()
				if tunnel.IsAllowRestart(err, true) {
					log.Info("tunnel restart in 10s...")
					time.Sleep(time.Second * 10)
					ch <- serv.Start()
				} else {
					app.Running = false
					return
				}
			} else {
				app.Running = false
				log.Info("tunnel exited")
				return
			}
		}
	}
}

//
// Run
// @Description: run application
// @receiver app
//
func (app *Application) Run() {
	var serv = tunnel.NewClient()
	if serv == nil {
		_ = log.Warn("tunnel server type not support")
		os.Exit(-1)
		return
	}
	app.serv = serv
	go app.runService(serv)
}

//
// Stop
// @Description:
// @receiver app
//
func (app *Application) Stop() {
	app.serv.Terminate()
}
