package main

import (
	"errors"
	"fmt"

	"github.com/sumit-behera-in/goLogger"
	"github.com/sumit-behera-in/gonas/p2p"
)

var logger *goLogger.Logger

func init() {
	var err error
	logger, err = goLogger.NewLogger("gonas", "log.log")
	if err != nil {
		panic("Failed to create logger instance : " + err.Error())
	}
}

func main() {
	logger.Info("Starting the application...")

	tcpOpts := p2p.TCPTransportOptions{
		ListenAddress: ":3000",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       &p2p.DefaultDecoder{},
		OnPeer:        onPeer,
		Logger:        logger,
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

func onPeer(p p2p.Peer) error {
	return errors.New("on peer failed")
}
