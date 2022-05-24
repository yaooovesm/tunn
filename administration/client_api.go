package administration

import (
	"errors"
	"github.com/gin-gonic/gin"
	"tunn/application"
	"tunn/config"
	"tunn/version"
)

//
// ApiCurrentApplication
// @Description:
// @param ctx
//
func ApiCurrentApplication(ctx *gin.Context) {
	running := false
	err := ""
	if application.Current != nil {
		running = application.Current.Running
		err = application.Current.Error
	}
	responseSuccess(ctx, map[string]interface{}{
		"running": running,
		"error":   err,
		"version": version.Version,
		"develop": version.Develop,
		"runtime": config.Current.Runtime,
	}, "")
}

//
// ApiApplicationStart
// @Description:
// @param ctx
//
func ApiApplicationStart(ctx *gin.Context) {
	if application.Current != nil && application.Current.Running {
		responseError(ctx, errors.New("客户端运行中"), "")
		return
	}
	cfg := config.ClientConfigStorage{}
	err := ctx.BindJSON(&cfg)
	if err != nil {
		response400(ctx)
		return
	}
	config.Current.Auth = config.Auth{
		Address: cfg.Auth.Address,
		Port:    cfg.Auth.Port,
	}
	config.Current.Security = config.Security{
		CertPem: cfg.Security.CertPem,
	}
	if cfg.User.Account != "" {
		config.Current.User.Account = cfg.User.Account
		config.Current.User.Password = cfg.User.Password
	}
	err = config.Current.Storage(true)
	if application.Current == nil {
		application.New()
	}
	application.Current.Run()
	responseSuccess(ctx, "", "启动成功")
}

//
// ApiApplicationStop
// @Description:
// @param ctx
//
func ApiApplicationStop(ctx *gin.Context) {
	if application.Current == nil || !application.Current.Running {
		responseError(ctx, errors.New("没有运行中的客户端"), "")
		return
	}
	application.Current.Stop()
	responseSuccess(ctx, "", "停止成功")

}
