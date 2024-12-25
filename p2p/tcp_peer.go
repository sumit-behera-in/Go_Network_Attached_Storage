package p2p

import "net"

// TCPPeer represents the remote node over an established TCP connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer.
	conn net.Conn

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
		conn:     conn,
		outbound: outbound,
	}
}

// Close closes the connection with the remote peer.
//
// Returns any error encountered while closing the connection.
func (t *TCPPeer) Close() error {
	return t.conn.Close()
}

// RemoteAddress retrieves the remote address of the peer.
//
// Returns the net.Addr of the remote peer (i.e., the address of the other side of the connection).
func (t *TCPPeer) RemoteAddress() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *TCPPeer) Send(data []byte) error {
	_, err := t.conn.Write(data)
	return err
}
