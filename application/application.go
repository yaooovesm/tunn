package application

import (
	log "github.com/cihub/seelog"
	"os"
	"sync"
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
	Config         config.Config
	Protocol       protocol.Name
	Serv           Service
	Running        bool
	Init           bool
	Error          string
	startWaitGroup *sync.WaitGroup
}

//
// New
// @Description:
// @return *Application
//
func New() *Application {
	app := &Application{
		Config:         config.Current,
		Protocol:       config.Current.Global.Protocol,
		Running:        false,
		Init:           false,
		startWaitGroup: &sync.WaitGroup{},
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
	Start(wg *sync.WaitGroup) error
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

func (app *Application) InitService() {
	log.Info("tunnel version : ", version.Version)
	if version.Develop {
		_ = log.Warn("当前版本为测试版本！")
	}
	//app.PProf()
	if app.Serv == nil {
		_ = log.Warn("tunnel server type not support")
		os.Exit(-1)
		return
	}
	if err := app.Serv.Init(); err != nil {
		_ = log.Warn("init failed : ", err)
		os.Exit(-1)
		return
	} else {
		log.Info("service init success...")
	}
	app.Init = true
}

//
// RunService
// @Description:
// @receiver app
// @param Serv
//
func (app *Application) RunService() {
	app.Error = ""
	ch := make(chan error, 1)
	go func() {
		app.Running = true
		ch <- app.Serv.Start(app.startWaitGroup)
	}()
	for {
		select {
		case err := <-ch:
			if err != nil {
				app.Serv.Stop()
				_ = log.Warn("tunnel exit with error : ", err.Error())
				app.Error = err.Error()
				if tunnel.IsAllowRestart(err, true) {
					log.Info("tunnel restart in 10s...")
					time.Sleep(time.Second * 10)
					app.startWaitGroup.Add(1)
					app.Error = ""
					ch <- app.Serv.Start(app.startWaitGroup)
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
func (app *Application) Run() error {
	if !app.Init {
		var serv = tunnel.NewClient()
		if serv == nil {
			return log.Warn("tunnel server type not support")
		}
		app.Serv = serv
		app.InitService()
	}
	app.startWaitGroup.Add(1)
	go app.RunService()
	app.startWaitGroup.Wait()
	return nil
}

//
// Stop
// @Description:
// @receiver app
//
func (app *Application) Stop() {
	app.Serv.Terminate()
}
