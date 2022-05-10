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

//
// Application
// @Description:
//
type Application struct {
	Config   config.Config
	Protocol protocol.Name
}

//
// New
// @Description:
// @return *Application
//
func New() *Application {
	return &Application{
		Config:   config.Current,
		Protocol: config.Current.Global.Protocol,
	}
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
}

//
// runService
// @Description:
// @receiver app
// @param serv
//
func (app *Application) runService(serv Service) {
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
	if err := serv.Init(); err != nil {
		_ = log.Warn("init failed : ", err)
		os.Exit(-1)
		return
	} else {
		log.Info("service init success...")
	}
	go func() {
		ch <- serv.Start()
	}()
	for {
		select {
		case err := <-ch:
			if err != nil {
				serv.Stop()
				_ = log.Warn("tunnel exit with error : ", err.Error())
				//if strings.Contains(str, "use of closed network connection")
				if tunnel.IsAllowRestart(err, true) {
					log.Info("tunnel restart in 10s...")
					time.Sleep(time.Second * 10)
					ch <- serv.Start()
				} else {
					os.Exit(-1)
					return
				}
			} else {
				log.Info("tunnel exited")
				os.Exit(0)
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
	app.runService(serv)
}
