package p2p

import (
	"net"
)

// always add mutex above the thing you want to protect

type TCPTransport struct {
	TCPTransportOptions
	listener     net.Listener
	responseChan chan Response
}

func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
		responseChan:        make(chan Response),
	}
}

// Dial implements the Transport interface, It is used to create outbound connections
func (t *TCPTransport) Dial(address string) error {
	con, err := net.Dial("tcp", address)
	if err != nil {
		t.Logger.Errorf("TCP Dial failed for address : %s", address)
		return err
	}

	go t.handleConn(con, true)

	return nil
}

// ListenAndAccept function is used to initialize the listener and accept
func (t *TCPTransport) ListenAndAccept() error {

	t.Logger.Infof("Initiating TCP to listen on %s ", t.ListenAddress)

	var err error

	// initialize the listener
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	t.Logger.Infof("TCP listen to  %s successful", t.ListenAddress)

	go t.startAcceptLoop()

	return nil

}

// accept connections asynchronously in a infinite loop
func (t *TCPTransport) startAcceptLoop() {

	t.Logger.Info("Starting TCP accept loop")

	for {
		// accept from the listener
		conn, err := t.listener.Accept()
		if err == net.ErrClosed {
			return
		}
		if err != nil {
			t.Logger.Errorf("Tcp accept error: %s", err)
		}

		t.Logger.Infof("Accepted TCP connection from %s", conn.RemoteAddr())

		go t.handleConn(conn, false)
	}
}

// handle the established connection

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	// create a new tcp peer
	peer := NewTCPPeer(conn, outbound)
	// use %+v fo more info on the parameters
	t.Logger.Infof("New incoming TCP connection %+v", peer)

	defer func() {
		t.Logger.Errorf("Dropping peer connection with error: %s", err.Error())
		conn.Close()
	}()

	if err = t.HandShakeFunc(peer); err != nil {
		t.Logger.Errorf("Handshake using TCP failed: %s", err)
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			t.Logger.Errorf("TCP OnPeer failed: %s, Maybe the TCP is connected with other peers", err)
			return
		}
	}

	// Read loop
	rpc := Response{}
	for {

		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			t.Logger.Errorf("TCP failed to decode payload from %+v : %s", conn, err)
			return
		}

		rpc.From = conn.RemoteAddr()
		t.responseChan <- rpc
		t.Logger.Infof("Response: %+v", rpc)
	}
}

// Consume implements transporter interface, which will return read only channel for reading the incoming messages received from another peer
func (t *TCPTransport) Consume() <-chan Response {
	// <- is used make read only channel
	return t.responseChan
}

// Close , closes the tcp listener
func (t *TCPTransport) Close() error {
	t.Logger.Infof("Dropping TCP connection with %s", t.ListenAddress)
	t.Logger.Close()
	return t.listener.Close()
}
