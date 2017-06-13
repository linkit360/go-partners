package src

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"runtime"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/linkit360/go-partners/server/src/config"
	"github.com/linkit360/go-partners/server/src/handlers"
	"github.com/linkit360/go-partners/service"
	m "github.com/linkit360/go-utils/metrics"
)

func Run() {
	appConfig := config.LoadConfig()

	service.Init(
		appConfig.AppName,
		appConfig.Service,
		appConfig.Notifier,
		appConfig.Mid,
	)

	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	log.WithField("CPUCount", nuCPU)

	go runGin(appConfig)
	runRPC(appConfig)
}

func runGin(appConfig config.AppConfig) {
	r := gin.New()
	m.AddHandler(r)

	r.Run(":" + appConfig.Server.HTTPPort)
	log.WithField("port", appConfig.Server.HTTPPort).Info("service port")
}

func runRPC(appConfig config.AppConfig) {

	l, err := net.Listen("tcp", "127.0.0.1:"+appConfig.Server.RPCPort)
	if err != nil {
		log.Fatal("netListen ", err.Error())
	} else {
		log.WithField("port", appConfig.Server.RPCPort).Info("rpc port")
	}

	server := rpc.NewServer()
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	server.RegisterName("Destination", &handlers.Destination{})

	for {
		if conn, err := l.Accept(); err == nil {
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		} else {
			log.WithField("error", err.Error()).Error("accept")
		}
	}
}
