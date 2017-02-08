package service

import (
	"fmt"

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
}

type GetDestinationParams struct {
	CountryCode  int64 `json:"country_code"`
	OperatorCode int64 `json:"operator_code"`
}

func GetDestination(p GetDestinationParams) (d inmem_service.Destination, err error) {
	sCount, err := inmem_client.GetAllRedirectStatCounts()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("cannot get all stats")
		return
	}
	for _, d = range svc.dsts {
		stat, ok := sCount[d.DestinationId]
		if !ok {
			log.WithFields(log.Fields{
				"dest_id": d.DestinationId,
			}).Error("not found in stat")
			continue
		}
		if stat.Count < d.AmountLimit && p.CountryCode == d.CountryCode {
			log.WithFields(log.Fields{
				"dest_id":    d.DestinationId,
				"partner_id": d.PartnerId,
				"url":        d.Destination,
			}).Debug("choose url")
			return d, nil
		}
	}
	err = fmt.Errorf("Not found: %d", p.CountryCode)
	return
}
