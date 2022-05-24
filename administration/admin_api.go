package administration

import (
	"github.com/gin-gonic/gin"
	"tunn/config"
	"tunn/version"
)

//
// ApiAdminInfo
// @Description:
// @param ctx
//
func ApiAdminInfo(ctx *gin.Context) {
	requireLogin := false
	if config.Current.Admin.Password != "" {
		requireLogin = true
	}
	responseSuccess(ctx, map[string]interface{}{
		"require_login": requireLogin,
		"version":       version.Version,
		"develop":       version.Develop,
	}, "")
}

//
// ApiAdminLogin
// @Description:
// @param ctx
//
func ApiAdminLogin(ctx *gin.Context) {
	if config.Current.Admin.Password == "" {
		responseSuccess(ctx, "", "登录成功")
		return
	}
	user := config.AdminUser{}
	err := ctx.BindJSON(&user)
	if err != nil {
		response400(ctx)
		return
	}
	err = config.Current.Admin.Check(user)
	if err != nil {
		responseError(ctx, err, "登录失败")
		return
	}
	responseSuccess(ctx, "", "登录成功")
}
