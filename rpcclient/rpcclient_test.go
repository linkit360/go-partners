package rpcclient

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	mid_service "github.com/linkit360/go-mid/service"
	partners_service "github.com/linkit360/go-partners/service"
)

func init() {
	c := RPCClientConfig{
		DSN:     "localhost:50312",
		Timeout: 10,
	}
	if err := Init(c); err != nil {
		log.WithField("error", err.Error()).Fatal("cannot init client")
	}
}

func TestGetUrlByRec(t *testing.T) {
	res, err := GetDestination(partners_service.GetDestinationParams{
		CountryCode:  66,
		OperatorCode: 515,
	})
	assert.Error(t, err, "error on unknown country")

	res, err = GetDestination(partners_service.GetDestinationParams{
		CountryCode:  92,
		OperatorCode: 515,
	})
	assert.Nil(t, err, "no error on ok country")
	//fmt.Printf("%#v %#v", res, err)

	expected := mid_service.Destination{
		DestinationId: 1,
		PartnerId:     1,
		AmountLimit:   0x3,
		Destination:   "http://default",
		RateLimit:     1,
		PricePerHit:   1,
		Score:         1,
		CountryCode:   92,
		OperatorCode:  41001,
	}
	assert.Equal(t, expected, res, "default url on empry record")
}
