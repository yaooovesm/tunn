package administration

import (
	"embed"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
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
	static embed.FS
}

//
// NewClientAdmin
// @Description:
// @param cfg
// @return *ClientAdmin
//
func NewClientAdmin(cfg config.Admin, static embed.FS) *ClientAdmin {
	if !version.Develop {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	return &ClientAdmin{
		cfg:    cfg,
		engine: engine,
		static: static,
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
	//setup static
	ad.engine.GET("/static/*filepath", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(ad.static))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
	//redirect
	ad.engine.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/static/"
		c.Redirect(http.StatusMovedPermanently, c.Request.URL.String())
	})
	//setup api
	api := NewClientAdminApi(ad.engine)
	api.Serv()
	log.Info("admin work at : ", address)
	accessLink := strings.Join([]string{ad.cfg.Address, strconv.Itoa(ad.cfg.Port)}, ":")
	if ad.cfg.Address == "0.0.0.0" {
		accessLink = strings.Join([]string{"127.0.0.1", strconv.Itoa(ad.cfg.Port)}, ":")
	}
	fmt.Println("for control please access link: ")
	fmt.Println("http://" + accessLink)
	err := ad.engine.Run(address)
	if err != nil {
		_ = log.Error("admin service stopped!")
		return
	}
}
