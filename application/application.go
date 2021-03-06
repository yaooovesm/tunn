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
	Serv           *tunnel.Client
	Running        bool
	Init           bool
	Error          string
	startWaitGroup *sync.WaitGroup
	terminate      chan int
	ch             chan error
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
		terminate:      make(chan int),
		ch:             make(chan error, 1),
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
Loop:
	for {
		select {
		case <-app.terminate:
			app.Error = "terminated"
			break Loop
		default:
			log.Info("tunnel start...")
			go func() {
				app.Running = true
				app.Error = ""
				app.ch <- app.Serv.Start(app.startWaitGroup)
			}()
			select {
			case err := <-app.ch:
				if err != nil {
					app.Serv.Stop()
					_ = log.Warn("tunnel exit with error : ", err.Error())
					app.Error = err.Error()
					if tunnel.IsAllowRestart(err, true) {
						log.Info("tunnel restart in 10s...")
						time.Sleep(time.Second * 10)
						goto Loop
						//app.Error = ""
						//app.ch <- app.Serv.Start(app.startWaitGroup)
					} else {
						break Loop
					}
				} else {
					break Loop
				}
			}
		}
	}
	app.Running = false
	log.Info("tunnel exited")
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
	if !app.Serv.Online {
		app.terminate <- 1
	} else {
		app.Serv.Terminate()
	}
	//app.ch <- errors.New("terminated")
}
