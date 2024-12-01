package main

import (
	"errors"
	"fmt"

	"github.com/sumit-behera-in/Distributed_File_Storage_GO/p2p"
)

func onPeer(p p2p.Peer) error {
	return errors.New("on peer failed")
}

func main() {
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddress: ":3000",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       &p2p.DefaultDecoder{},
		OnPeer:        onPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {

		for {
			println("data using channel :")
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
			println(":  data using channel ")
		}

	}()

	if err := tr.ListenAndAccept(); err != nil {
		panic(err.Error())
	}

	fmt.Println("HEllo")

	select {}
}
