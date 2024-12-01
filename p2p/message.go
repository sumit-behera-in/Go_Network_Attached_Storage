package p2p

// Message hold arbitrary data that is being sent over each transport between two nodes in a network.
type Message struct {
	Payload []byte
}
