package service

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"

	inmem_client "github.com/linkit360/go-inmem/rpcclient"
	inmem_service "github.com/linkit360/go-inmem/service"
	"github.com/linkit360/go-partners/server/src/config"
	"github.com/linkit360/go-utils/amqp"
	"github.com/linkit360/go-utils/cqr"
	m "github.com/linkit360/go-utils/metrics"
)

var svc Service

type Service struct {
	conf      config.ServiceConfig
	dsts      []inmem_service.Destination
	cqrConfig []cqr.CQRConfig
	m         *Metrics
	db        *sql.DB
	notifier  *amqp.Notifier
}

type Metrics struct {
	Success                m.Gauge
	Errors                 m.Gauge
	GetDestinationDuration prometheus.Summary
}

type DestinationHit struct {
	DestinationHitId int64     `json:"id"`
	PartnerId        int64     `json:"partner_id"`
	DestinationId    int64     `json:"destination_id"`
	PricePerHit      float64   `json:"price_per_hit,omitempty"`
	Tid              string    `json:"tid"`
	CreatedAt        time.Time `json:"cerated_at"`
	SentAt           time.Time `json:"sent_at"`
	Destination      string    `json:"destination"`
	Msisdn           string    `json:"msisdn"`
	OperatorCode     int64     `json:"operator_code"`
	CountryCode      int64     `json:"country_code"`
	ResponseCode     int       `json:"response_code"`
}

func Init(
	appName string,
	serviceConfig config.ServiceConfig,
	notifierConfig amqp.NotifierConfig,
	inMemConfig inmem_client.ClientConfig,
) {
	log.SetLevel(log.DebugLevel)
	svc.conf = serviceConfig
	svc.notifier = amqp.NewNotifier(notifierConfig)
	svc.m = initMetrics(appName)

	if err := inmem_client.Init(inMemConfig); err != nil {
		log.WithField("error", err.Error()).Fatal("cannot init inmem client")
	}

	// reload
	reloadDestinations()
	go func() {
		for range time.Tick(time.Duration(10) * time.Minute) {
			reloadDestinations()
		}
	}()

}
func initMetrics(appName string) *Metrics {

	appM := &Metrics{
		Success:                m.NewGauge("", "", "success", "success"),
		Errors:                 m.NewGauge("", "", "errors", "errors"),
		GetDestinationDuration: m.NewSummary(appName+"_get_destination_duration_seconds", "hit duration seconds"),
	}

	go func() {
		for range time.Tick(time.Minute) {
			appM.Success.Update()
			appM.Errors.Update()
		}
	}()
	return appM
}
