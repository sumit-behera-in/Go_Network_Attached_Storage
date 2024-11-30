package p2p

import "net"

// TCPPeer represent the remote node over a established TCP connection
type TCPPeer struct {
	// it is the underlying connection of the peer
	conn net.Conn

	// if we dial and retrieve a connection => outbound = true
	// if we accept and retrieve a connection => outbound = false (inbound)
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}
