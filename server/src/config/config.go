package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/configor"

	inmem_client "github.com/vostrok/inmem/rpcclient"
	"github.com/vostrok/utils/amqp"
	"github.com/vostrok/utils/db"
)

type AppConfig struct {
	AppName   string                       `yaml:"app_name"`
	Server    ServerConfig                 `yaml:"server"`
	Service   ServiceConfig                `yaml:"service"`
	InMem     inmem_client.RPCClientConfig `yaml:"inmem"`
	DB        db.DataBaseConfig            `yaml:"db"`
	Publisher amqp.NotifierConfig          `yaml:"notifier"`
}

type ServiceConfig struct {
	Queues          QueuesConfig `yaml:"queues"`
	ResponseLogPath string       `default:"/var/log/linkit/partner_requests.log" yaml:"response"`
}
type QueuesConfig struct {
	NewHitNotify string `yaml:"hits"`
}
type ServerConfig struct {
	Port string `default:"50311"`
}

func LoadConfig() AppConfig {
	cfg := flag.String("config", "dev/partners.yml", "configuration yml file")
	flag.Parse()
	var appConfig AppConfig

	if *cfg != "" {
		if err := configor.Load(&appConfig, *cfg); err != nil {
			log.WithField("config", err.Error()).Fatal("config load error")
		}
	}

	if appConfig.AppName == "" {
		log.Fatal("app_name must be defiled as <host>_<name>")
	}
	if strings.Contains(appConfig.AppName, "-") {
		log.Fatal("app_name must be without '-' : it's not a valid metric name")
	}
	appConfig.Server.Port = envString("PORT", appConfig.Server.Port)
	appConfig.Publisher.Conn.Host = envString("RBMQ_HOST", appConfig.Publisher.Conn.Host)

	log.WithField("config", fmt.Sprintf("%#v", appConfig)).Info("Config loaded")
	return appConfig
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
