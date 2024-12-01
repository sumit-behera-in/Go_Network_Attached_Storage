package p2p

import "net"

// Response hold arbitrary data that is being sent over each transport between two nodes in a network.
type Response struct {
	From    net.Addr
	Payload []byte
}
