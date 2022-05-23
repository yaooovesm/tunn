package administration

import (
	"github.com/gin-gonic/gin"
	"tunn/config"
)

//
// ApiSaveCurrentConfig
// @Description:
// @param ctx
//
func ApiSaveCurrentConfig(ctx *gin.Context) {
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
	}
	err = config.Current.Storage(true)
	if err != nil {
		responseError(ctx, err, "保存失败")
		return
	}
	responseSuccess(ctx, "", "保存成功")
}

//
// ApiGetCurrentConfig
// @Description:
// @param ctx
//
func ApiGetCurrentConfig(ctx *gin.Context) {
	cfg := config.Current.ToStorageModel()
	admin := config.Admin{
		Address: cfg.Admin.Address,
		Port:    cfg.Admin.Port,
	}
	cfg.Admin = admin
	responseSuccess(ctx, cfg, "")
}
