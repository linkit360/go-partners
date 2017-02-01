package handlers

import (
	"github.com/vostrok/partners/src/service"
	"github.com/vostrok/utils/rec"
)

type GetByRecordParams struct {
	Record rec.Record `json:"record,omitempty"`
}

type Url struct{}

func (rpc *Url) GetByRecord(
	req GetByRecordParams, res *string) error {

	url, err := service.GetByRecord(req)
	if err != nil {
		return err
	}
	*res = url
	return nil
}
