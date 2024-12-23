package p2p

import "net"

// Response hold arbitrary data that is being sent over each transport between two nodes in a network.
type Response struct {
	From    net.Addr // From is the address of the node that sent the data
	Payload []byte   // Payload is the data that is being sent over the network
}
