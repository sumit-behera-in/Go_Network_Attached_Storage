package p2p

import (
	"fmt"
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

// ListenAndAccept function is used to initialize the listener and accept
func (t *TCPTransport) ListenAndAccept() error {

	var err error

	// initialize the listener
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil

}

// accept connections asynchronously in a infinite loop
func (t *TCPTransport) startAcceptLoop() {
	for {
		// accept from the listener
		conn, err := t.listener.Accept()

		if err != nil {
			fmt.Printf("Tcp accept error: %s\n", err)
		}

		go t.handleConn(conn)
	}
}

// handle the established connection

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	// create a new tcp peer
	peer := NewTCPPeer(conn, true)
	// use %+v fo more info on the parameters
	fmt.Printf("new incoming connection %+v\n", peer)

	defer func() {
		fmt.Printf("Dropping peer connection with error: %s\n", err.Error())
		conn.Close()
	}()

	if err = t.HandShakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := Response{}
	for {

		err = t.Decoder.Decode(conn, &rpc)

		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr()
		t.responseChan <- rpc
		fmt.Printf("Response: %+v\n", rpc)
	}
}

// Consume implements transporter interface, which will return read only channel for reading the incoming messages received from another peer
func (t *TCPTransport) Consume() <-chan Response {
	// <- is used make read only channel
	return t.responseChan
}
