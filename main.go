package main

import (
	"github.com/sumit-behera-in/goLogger"
	"github.com/sumit-behera-in/gonas/fileserver"
	"github.com/sumit-behera-in/gonas/p2p"
	"github.com/sumit-behera-in/gonas/storage"
)

var logger *goLogger.Logger

func init() {
	var err error
	logger, err = goLogger.NewLogger("gonas", "", 10, 4, "IST")
	if err != nil {
		panic("Failed to create logger instance : " + err.Error())
	}
}

func main() {
	logger.Info("Starting the application...")

	fileServer1 := makeServer(":1000", "NAS_1000_root", ":2000", ":3000", ":4000")
	fs2 := makeServer(":2000", "NAS_2000_root", ":1000", ":3000", ":4000")
	fs3 := makeServer(":3000", "NAS_2000_root", ":1000", ":2000", ":4000")
	fs4 := makeServer(":4000", "NAS_2000_root", ":1000", ":3000", ":2000")

	go func() {
		fileServer1.Start()
	}()

	go func() {
		fs2.Start()
	}()

	go func() {
		fs3.Start()
	}()

	go func() {
		fs4.Start()
	}()

	select {}

}

func makeServer(listenAddress string, storageRoot string, nodes ...string) *fileserver.Fileserver {
	logger.Infof("Creating a file server with address : %s", listenAddress)
	tcpTransportOpts := p2p.TCPTransportOptions{
		Logger:        logger,
		ListenAddress: listenAddress,
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       &p2p.GOBDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOptions := fileserver.FileServerOpts{
		Logger:            logger,
		StorageRoot:       storageRoot,
		PathTransformFunc: storage.CASPathTransformFunc,
		Transport:         tcpTransport,
		BootStrapNodes:    nodes,
	}
	server := fileserver.NewFileServer(fileServerOptions)
	tcpTransport.OnPeer = server.OnPeer

	return server

}
