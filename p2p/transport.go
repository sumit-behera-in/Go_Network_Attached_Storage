package p2p

// Peer is an interface that represents the remote node
type Peer interface {
}

// transport is anything that handles communication
// btn nodes in network, this can be (tcp , udp , websocket ...)
type Transport interface {
	ListenAndAccept() error
}
