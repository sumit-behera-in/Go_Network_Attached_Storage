package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransportOptions struct {
	ListenAddress string
	HandShakeFunc HandShakerFunc
	Decoder       Decoder
}

// always add mutex above the thing you want to protect

type TCPTransport struct {
	TCPTransportOptions
	listener     net.Listener
	responseChan chan Response

	mu    sync.RWMutex
	peers map[net.Addr]Peer
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
	// create a new tcp peer
	peer := NewTCPPeer(conn, true)
	// use %+v fo more info on the parameters
	fmt.Printf("new incoming connection %+v\n", peer)

	if err := t.HandShakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("Hand Shaking failed for connection: %+v \n", peer)
		return
	}

	// Read loop
	rpc := Response{}
	for {

		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP error: %s\n", err.Error())
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
