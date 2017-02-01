package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"

	inmem_client "github.com/vostrok/inmem/rpcclient"
	"github.com/vostrok/partners/server/src/config"
	"github.com/vostrok/utils/amqp"
	m "github.com/vostrok/utils/metrics"
	"github.com/vostrok/utils/rec"
)

var svc Service

type EventNotify struct {
	EventName string          `json:"event_name,omitempty"`
	EventData NewHitNotifyMsg `json:"event_data,omitempty"`
}

type Service struct {
	notifier *amqp.Notifier
	db       *sql.DB
	conf     config.ServiceConfig
	m        *Metrics
}

type Metrics struct {
	Success     m.Gauge
	Errors      m.Gauge
	HitDuration prometheus.Summary
}

type NewHitNotifyMsg struct {
	R            rec.Record `yaml:"record"`
	Referer      string     `yaml:"referer"`
	Url          string     `yaml:"url"`
	ResponseCode int        `yaml:"response_code"`
	HitAt        time.Time  `yaml:"hit_at"`
}

func InitService(
	appName string,
	serviceConfig config.ServiceConfig,
	notifierConfig amqp.NotifierConfig,
	inMemConfig inmem_client.RPCClientConfig,
) {
	log.SetLevel(log.DebugLevel)
	svc.conf = serviceConfig
	svc.notifier = amqp.NewNotifier(notifierConfig)
	svc.m = initMetrics(appName)

	if err := inmem_client.Init(inMemConfig); err != nil {
		log.WithField("error", err.Error()).Fatal("cannot init inmem client")
	}
}
func initMetrics(appName string) *Metrics {

	appM := &Metrics{
		Success:     m.NewGauge("", "", "success", "success"),
		Errors:      m.NewGauge("", "", "errors", "errors"),
		HitDuration: m.NewSummary(appName+"_hit_duration_seconds", "hit duration seconds"),
	}

	go func() {
		for range time.Tick(time.Minute) {
			appM.Success.Update()
			appM.Errors.Update()
		}
	}()
	return appM
}

func (svc *Service) newHitNotify(msg NewHitNotifyMsg) error {
	if msg.HitAt.IsZero() {
		msg.HitAt = time.Now().UTC()
	}

	event := EventNotify{
		EventName: svc.conf.Queues.NewHitNotify,
		EventData: msg,
	}
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json.Marshal: %s", err.Error())
	}
	svc.notifier.Publish(amqp.AMQPMessage{svc.conf.Queues.NewHitNotify, uint8(1), body})
	return nil
}
