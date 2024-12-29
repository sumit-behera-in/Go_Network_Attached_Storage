package p2p

import "net"

// TCPPeer represents the remote node over an established TCP connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer, which is tcp connection.
	net.Conn

	// outbound indicates if the connection is outgoing (true) or incoming (false).
	outbound bool
}

// NewTCPPeer creates a new TCPPeer instance.
//
// conn: The TCP connection object that represents the connection to the remote peer.
// outbound: A boolean flag indicating whether this connection is outgoing (true) or incoming (false).
//
// Returns a new instance of TCPPeer.
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

func (t *TCPPeer) Send(data []byte) error {
	_, err := t.Conn.Write(data)
	return err
}
