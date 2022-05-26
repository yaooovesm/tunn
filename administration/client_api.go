package administration

import (
	"errors"
	"github.com/gin-gonic/gin"
	"tunn/application"
	"tunn/authenticationv2"
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
	initialized := false
	online := false
	err := ""
	if application.Current != nil {
		running = application.Current.Running
		err = application.Current.Error
		initialized = application.Current.Init
		if application.Current.Serv != nil {
			online = application.Current.Serv.Online
		}
	}
	responseSuccess(ctx, map[string]interface{}{
		"running":     running,
		"online":      online,
		"initialized": initialized,
		"error":       err,
		"version":     version.Version,
		"develop":     version.Develop,
		"runtime":     config.Current.Runtime,
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
	err = application.Current.Run()
	if err != nil {
		responseError(ctx, err, "启动失败")
		return
	}
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

//
// ApiGetAvailableExports
// @Description:
// @param ctx
//
func ApiGetAvailableExports(ctx *gin.Context) {
	if application.Current == nil || !application.Current.Running {
		responseError(ctx, errors.New("没有运行中的客户端"), "")
		return
	}
	res, err := application.Current.Serv.AuthClient.Operation(authenticationv2.OperationGetAvailableExports, nil)
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	responseSuccess(ctx, res.Result, "")
}

//
// ApiGetConfig
// @Description:
// @param ctx
//
func ApiGetConfig(ctx *gin.Context) {
	if application.Current == nil || !application.Current.Running {
		responseError(ctx, errors.New("没有运行中的客户端"), "")
		return
	}
	res, err := application.Current.Serv.AuthClient.Operation(authenticationv2.OperationGetConfig, nil)
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	responseSuccess(ctx, res.Result, "")
}

//
// ApiUpdateRoutes
// @Description:
// @param ctx
//
func ApiUpdateRoutes(ctx *gin.Context) {
	if application.Current == nil || !application.Current.Running {
		responseError(ctx, errors.New("没有运行中的客户端"), "")
		return
	}
	var routes []config.Route
	err := ctx.BindJSON(&routes)
	if err != nil {
		response400(ctx)
		return
	}
	result, err := application.Current.Serv.AuthClient.Operation(authenticationv2.OperationUpdateRoutes, map[string]interface{}{
		"routes": routes,
	})
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	if result.Error == "" {
		responseSuccess(ctx, "", "保存成功")
	} else {
		responseError(ctx, errors.New(result.Error), "")
	}
}

//
// ApiResetRoutes
// @Description:
// @param ctx
//
func ApiResetRoutes(ctx *gin.Context) {
	if application.Current == nil || !application.Current.Running {
		responseError(ctx, errors.New("没有运行中的客户端"), "")
		return
	}
	result, err := application.Current.Serv.AuthClient.Operation(authenticationv2.OperationResetRoutes, nil)
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	if result.Error == "" {
		responseSuccess(ctx, "", "保存成功")
	} else {
		responseError(ctx, errors.New(result.Error), "")
	}
}
