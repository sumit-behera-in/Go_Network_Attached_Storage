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
	handShakeFunc    HandShakerFunc

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddress string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddress,
		handShakeFunc: func(a any) error {return nil},
	}
}

// initalize the listener and accept
func (t *TCPTransport) ListenAndAccept() error {

	var err error

	// initalize the listener
	t.listener, err = net.Listen("tcp", t.listenAddress)

	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil

}

// accept connetions asynchronously in a infinite loop

func (t *TCPTransport) startAcceptLoop() {
	for {
		// accept from the listner
		conn, err := t.listener.Accept()

		if err != nil {
			fmt.Printf("Tcp accept error: %v\n", err)
		}

		go t.handleConn(conn)
	}
}

// handle the established connection

func (t *TCPTransport) handleConn(conn net.Conn) {
	// create a new tcp peer
	peer := NewTCPPeer(conn, true)

	fmt.Println("new incoming connection", peer)
}
