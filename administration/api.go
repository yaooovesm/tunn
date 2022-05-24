package administration

import "github.com/gin-gonic/gin"

//
// ClientAdminApi
// @Description:
//
type ClientAdminApi struct {
	r *gin.Engine
}

//
// NewClientAdminApi
// @Description:
// @param engine
// @return *ClientAdminApi
//
func NewClientAdminApi(r *gin.Engine) *ClientAdminApi {
	return &ClientAdminApi{r: r}
}

//
// Serv
// @Description:
// @receiver api
//
func (ad *ClientAdminApi) Serv() {
	api := ad.r.Group("/api")
	//config
	{
		configGroup := api.Group("/config")
		configGroup.GET("/", ApiGetCurrentConfig)
		configGroup.GET("/all", ApiGetCurrentConfigAll)
		configGroup.POST("/save", ApiSaveCurrentConfig)
	}
	//application
	{
		appGroup := api.Group("/application")
		appGroup.GET("/", ApiCurrentApplication)
		appGroup.POST("/start", ApiApplicationStart)
		appGroup.GET("/stop", ApiApplicationStop)
	}
	//admin
	{
		adminGroup := api.Group("/admin")
		adminGroup.GET("/", ApiAdminInfo)
		adminGroup.POST("/login", ApiAdminLogin)
	}
}
