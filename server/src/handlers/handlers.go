package handlers

import (
	"github.com/vostrok/partners/service"
)

type Url struct{}

func (rpc *Url) GetByRecord(
	req service.GetDestinationParams, res *string) error {

	url, err := service.GetDestination(req)
	if err != nil {
		return err
	}
	*res = url
	return nil
}
