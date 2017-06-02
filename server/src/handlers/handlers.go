package handlers

import (
	mid_service "github.com/linkit360/go-mid/service"
	"github.com/linkit360/go-partners/service"
)

type Destination struct{}

func (rpc *Destination) Get(
	req service.GetDestinationParams, res *mid_service.Destination) error {

	dst, err := service.GetDestination(req)
	if err != nil {
		return err
	}
	*res = dst
	return nil
}
