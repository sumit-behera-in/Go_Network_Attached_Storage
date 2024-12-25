package p2p

import (
	"net"
)

// TCPTransport represents the TCP transport layer for handling peer-to-peer communication.
type TCPTransport struct {
	TCPTransportOptions
	listener     net.Listener    // Listener for accepting incoming TCP connections.
	responseChan chan Response   // Channel for receiving responses from remote peers.
}

// NewTCPTransport creates a new TCPTransport instance with the given options.
//
// opts: The TCP transport options that configure the transport behavior.
//
// Returns a pointer to a new TCPTransport instance.
func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
		responseChan:        make(chan Response),
	}
}

// Dial establishes an outbound TCP connection to the given address and starts handling the connection.
//
// address: The address of the remote peer to connect to.
//
// Returns an error if the dial operation fails.
func (t *TCPTransport) Dial(address string) error {
	con, err := net.Dial("tcp", address)
	if err != nil {
		t.Logger.Errorf("TCP Dial failed for address : %s", address)
		return err
	}

	// Handle the established connection in a goroutine.
	go t.handleConn(con, true)

	return nil
}

// ListenAndAccept starts listening for incoming TCP connections on the specified address.
//
// Returns an error if there is an issue initializing the listener or starting the accept loop.
func (t *TCPTransport) ListenAndAccept() error {
	t.Logger.Infof("Initiating TCP to listen on %s ", t.ListenAddress)

	// Initialize the listener.
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	t.Logger.Infof("TCP listen to %s successful", t.ListenAddress)

	// Start the loop for accepting incoming connections asynchronously.
	go t.startAcceptLoop()

	return nil
}

// startAcceptLoop is an infinite loop that accepts incoming TCP connections from the listener.
//
// It handles each new connection in a separate goroutine.
func (t *TCPTransport) startAcceptLoop() {
	t.Logger.Info("Starting TCP accept loop")

	for {
		// Accept a connection from the listener.
		conn, err := t.listener.Accept()
		if err == net.ErrClosed {
			return
		}
		if err != nil {
			t.Logger.Errorf("Tcp accept error: %s", err)
		}

		t.Logger.Infof("Accepted TCP connection from %s", conn.RemoteAddr())

		// Handle the accepted connection in a separate goroutine.
		go t.handleConn(conn, false)
	}
}

// handleConn handles the established TCP connection, performs the handshake, and processes incoming data.
//
// conn: The TCP connection to the remote peer.
// outbound: A boolean flag indicating whether this connection is outbound (true) or inbound (false).
func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	// Create a new TCPPeer instance to represent the connected peer.
	peer := NewTCPPeer(conn, outbound)
	t.Logger.Infof("New incoming TCP connection %+v", peer)

	// Ensure that the connection is properly closed on exit.
	defer func() {
		t.Logger.Errorf("Dropping peer connection with error: %s", err.Error())
		conn.Close()
	}()

	// Perform the handshake with the remote peer.
	if err = t.HandShakeFunc(peer); err != nil {
		t.Logger.Errorf("Handshake using TCP failed: %s", err)
		return
	}

	// If the OnPeer callback is provided, invoke it to handle the peer.
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			t.Logger.Errorf("TCP OnPeer failed: %s, Maybe the TCP is connected with other peers", err)
			return
		}
	}

	// Read loop: Continuously receive data from the connection.
	rpc := Response{}
	for {
		// Decode the incoming data into a Response object.
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			t.Logger.Errorf("TCP failed to decode payload from %+v : %s", conn, err)
			return
		}

		// Set the remote address for the response and send it to the response channel.
		rpc.From = conn.RemoteAddr()
		t.responseChan <- rpc
		t.Logger.Infof("Response: %+v", rpc)
	}
}

// Consume implements the Transport interface, returning a read-only channel to consume incoming responses.
//
// Returns a read-only channel (chan Response) that allows the caller to receive incoming messages from remote peers.
func (t *TCPTransport) Consume() <-chan Response {
	// The channel is read-only, indicated by <-chan.
	return t.responseChan
}

// Close closes the TCP listener and performs any necessary cleanup operations.
//
// Returns an error if there is an issue closing the listener.
func (t *TCPTransport) Close() error {
	t.Logger.Infof("Dropping TCP connection with %s", t.ListenAddress)
	t.Logger.Close()  // Close the logger
	return t.listener.Close() // Close the TCP listener
}
