package p2p

// interface of remote node
type Peer interface {
}

// tansport is anything that handles communication
// btn nodes in network, this can be tcp , udp , websockets ...
type Transport interface {
	ListenAndAccept() error
}
