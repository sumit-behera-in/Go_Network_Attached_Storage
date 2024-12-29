package p2p

import (
	"net"

	"github.com/sumit-behera-in/goLogger"
)

// Peer is an interface that represents the remote node
type Peer interface {
	net.Conn
	Send([]byte) error
}

// transport is anything that handles communication
// btn nodes in network, this can be (tcp , udp , websocket ...)
type Transport interface {
	Dial(string) error
	ListenAndAccept() error   // ListenAndAccept is the function that is used to start the transport
	Consume() <-chan Response // Consume is the function that is used to consume the incoming data
	Close() error             // Close is the function that is used to stop the transport
}

type TCPTransportOptions struct {
	Logger        *goLogger.Logger // Logger is the logger instance
	ListenAddress string           // ListenAddress is the address on which the transport listens
	HandShakeFunc HandShakerFunc   // HandShakeFunc is the function that is called when a new connection is established
	Decoder       Decoder          // Decoder is the decoder that is used to decode the incoming data
	OnPeer        func(Peer) error // OnPeer is the function that is called when a new peer is connected to check if the peer is connected to other peers
}
