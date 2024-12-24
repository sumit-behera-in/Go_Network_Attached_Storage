package main

import (
	"github.com/sumit-behera-in/goLogger"
	fileserver "github.com/sumit-behera-in/gonas/fileServer"
	"github.com/sumit-behera-in/gonas/p2p"
	"github.com/sumit-behera-in/gonas/storage"
)

var logger *goLogger.Logger

func init() {
	var err error
	logger, err = goLogger.NewLogger("gonas", "./log", 1, 4, "IST")
	if err != nil {
		panic("Failed to create logger instance : " + err.Error())
	}
}

func main() {
	logger.Info("Starting the application...")
	defer logger.Close()

	tcpTransportOpts := p2p.TCPTransportOptions{
		Logger:        logger,
		ListenAddress: ":1000",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       &p2p.GOBDecoder{},
		// TODO: onPeer func
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOptions := fileserver.FileServerOpts{
		Logger:            logger,
		StorageRoot:       "NAS_1000_root",
		PathTransformFunc: storage.CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	fileServer := fileserver.NewFileServer(fileServerOptions)

	if err := fileServer.Start(); err !=  nil {
		logger.Fatal(err.Error())
	}

	select {}
}
