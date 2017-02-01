package rpcclient

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vostrok/utils/rec"
)

func init() {
	c := RPCClientConfig{
		DSN:        "localhost:50311",
		Timeout:    10,
		DefaultUrl: "http://default",
	}
	if err := Init(c); err != nil {
		log.WithField("error", err.Error()).Fatal("cannot init client")
	}
}

func TestGetUrlByRec(t *testing.T) {
	r := rec.Record{
		Tid: rec.GenerateTID(),
	}
	res, err := GetUrlByRec(r)
	assert.Nil(t, err)
	expected := cli.conf.DefaultUrl
	assert.Equal(t, expected, res, "default url on empry record")
}
