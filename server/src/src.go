package src

import (
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/vostrok/partners/src/server/config"
	"github.com/vostrok/partners/src/service"
	m "github.com/vostrok/utils/metrics"
)

func RunServer() {
	appConfig := config.LoadConfig()

	service.InitService(
		appConfig.AppName,
		appConfig.Service,
		appConfig.Publisher,
		appConfig.InMem,
	)

	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	log.WithField("CPUCount", nuCPU)

	r := gin.New()
	m.AddHandler(r)

	r.NoRoute(notFound)

	r.Run(":" + appConfig.Server.Port)

	log.WithField("port", appConfig.Server.Port).Info("init")
}

func notFound(c *gin.Context) {
	log.WithFields(log.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"req":    c.Request.URL.RawQuery,
	}).Info("404notfound")
	c.JSON(200, struct{}{})
}
