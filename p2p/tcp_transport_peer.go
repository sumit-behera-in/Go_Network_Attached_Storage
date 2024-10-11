package p2p

import "net"

// represent the remote node over a established tcp connection
type TCPPeer struct {
	//it is the underlying connection of the peer
	conn net.Conn

	// if we dial and retrieve a connecton => outbound = true
	// if we accept and retrieve a connection => outound = false (inbound)
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}
