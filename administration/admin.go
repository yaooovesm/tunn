package administration

import (
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
	"strings"
	"tunn/config"
	"tunn/version"
)

//
// ClientAdmin
// @Description:
//
type ClientAdmin struct {
	cfg    config.Admin
	engine *gin.Engine
}

//
// NewClientAdmin
// @Description:
// @param cfg
// @return *ClientAdmin
//
func NewClientAdmin(cfg config.Admin) *ClientAdmin {
	if !version.Develop {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	return &ClientAdmin{
		cfg:    cfg,
		engine: engine,
	}
}

//
// Run
// @Description:
// @receiver ad
//
func (ad *ClientAdmin) Run() {
	if ad.cfg.Address == "" {
		ad.cfg.Address = "0.0.0.0"
	}
	if ad.cfg.Port <= 0 || ad.cfg.Port > 65535 {
		_ = log.Error("invalid admin port : ", ad.cfg.Port)
		return
	}
	var address string
	ip := net.ParseIP(ad.cfg.Address)
	if ip != nil {
		address = strings.Join([]string{ad.cfg.Address, strconv.Itoa(ad.cfg.Port)}, ":")
	} else {
		address = strings.Join([]string{"0.0.0.0", strconv.Itoa(ad.cfg.Port)}, ":")
	}
	//setup api
	api := NewClientAdminApi(ad.engine)
	api.Serv()
	log.Info("admin work at : ", address)
	err := ad.engine.Run(address)
	if err != nil {
		_ = log.Error("admin service stopped!")
		return
	}
}
