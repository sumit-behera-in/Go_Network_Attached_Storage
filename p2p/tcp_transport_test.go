package p2p

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOptions{
		ListenAddress: ":3000",
		HandShakeFunc: NOPHandShakeFunc,
		Decoder:       &DefaultDecoder{},
	}

	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddress, opts.ListenAddress)
	assert.Equal(t, tr.ListenAndAccept(), nil)
}
