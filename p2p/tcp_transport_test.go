package p2p

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/sumit-behera-in/goLogger"
)

func TestTCPTransport(t *testing.T) {
	logger, _ := goLogger.NewLogger("test", "", 1000, 2, "IST")
	opts := TCPTransportOptions{
		ListenAddress: ":3000",
		HandShakeFunc: NOPHandShakeFunc,
		Decoder:       &DefaultDecoder{},
		Logger:        logger,
	}

	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddress, opts.ListenAddress)
	assert.Equal(t, tr.ListenAndAccept(), nil)
}
