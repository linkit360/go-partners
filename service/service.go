package service

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	inmem_client "github.com/vostrok/inmem/rpcclient"
	inmem_service "github.com/vostrok/inmem/service"
)

func reloadDestinations() {
	dd, err := inmem_client.GetAllDestinations()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("cannot get all destinations")
	}
	svc.dsts = dd
	log.WithFields(log.Fields{
		"destinations": svc.dsts,
	}).Debug("destinations reloaded")
}

type GetDestinationParams struct {
	CountryCode  int64 `json:"country_code"`
	OperatorCode int64 `json:"operator_code"`
}

func GetDestination(p GetDestinationParams) (d inmem_service.Destination, err error) {
	begin := time.Now()
	log.WithFields(log.Fields{
		"country_code":  p.CountryCode,
		"operator_code": p.OperatorCode,
	}).Debug("go request")
	defer func() {
		svc.m.GetDestinationDuration.Observe(time.Since(begin).Seconds())
		if err == nil {
			svc.m.Success.Inc()
		} else {
			svc.m.Errors.Inc()
		}
	}()
	sCount, err := inmem_client.GetAllRedirectStatCounts()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("cannot get all stats")
		return
	}
	for _, d = range svc.dsts {
		log.WithFields(log.Fields{
			"dest_id":      d.DestinationId,
			"partner_id":   d.PartnerId,
			"country_code": d.CountryCode,
			"url":          d.Destination,
		}).Debug("considering..")

		// if not found in stat
		// then no stats yet and we can  add new
		stat, _ := sCount[d.DestinationId]
		if stat == nil {
			stat = &inmem_service.StatCount{
				DestinationId: d.DestinationId,
				Count:         0,
			}
		}
		log.WithFields(log.Fields{
			"dest_id": stat.DestinationId,
			"count":   stat.Count,
		}).Debug("stats")

		if stat.Count < d.AmountLimit && p.CountryCode == d.CountryCode {
			log.WithFields(log.Fields{
				"dest_id":    d.DestinationId,
				"partner_id": d.PartnerId,
				"url":        d.Destination,
			}).Info("choose url")
			return d, nil
		} else {
			log.WithFields(log.Fields{
				"country_check": (p.CountryCode == d.CountryCode),
				"limit_check":   (stat.Count < d.AmountLimit),
			}).Debug("not passed")
		}
	}
	err = fmt.Errorf("Destination for country %d not found", p.CountryCode)
	log.WithFields(log.Fields{
		"country_code": p.CountryCode,
	}).Error("not found")
	return
}
