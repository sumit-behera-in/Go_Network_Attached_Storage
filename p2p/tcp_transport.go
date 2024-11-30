package p2p

import (
	"fmt"
	"net"
	"sync"
)

// always add mutex above the thing you want to protect

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	handShakeFunc HandShakerFunc

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddress string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddress,
		handShakeFunc: NOPHandShakeFunc,
	}
}

// ListenAndAccept function is used to initialize the listener and accept
func (t *TCPTransport) ListenAndAccept() error {

	var err error

	// initialize the listener
	t.listener, err = net.Listen("tcp", t.listenAddress)
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
}
