package fileserver

import (
	"io"

	"github.com/sumit-behera-in/goLogger"
	"github.com/sumit-behera-in/gonas/p2p"
	"github.com/sumit-behera-in/gonas/storage"
)

type FileServerOpts struct {
	Logger            *goLogger.Logger          // this is the global logger
	StorageRoot       string                    // root storage to store all files or folders managed by goNAS
	PathTransformFunc storage.PathTransformFunc // used encrypt path and file name
	Transport         p2p.Transport             // TCP, UDP, HTTP
}

type Fileserver struct {
	FileServerOpts
	storage  storage.Storage
	quitChan chan struct{}
}

func NewFileServer(option FileServerOpts) *Fileserver {
	storageOpts := storage.StorageOptions{
		PathTransformFunc: option.PathTransformFunc,
	}

	return &Fileserver{
		FileServerOpts: option,
		storage:        *storage.NewStorage(storageOpts),
		quitChan:       make(chan struct{}),
	}
}

func (server *Fileserver) keepAlive() {
	server.Logger.Info("File Server KeepAlive() is called")

	defer func() {
		server.Logger.Info("File Server KeepAlive() is stopped")
		server.Transport.Close()
	}()

	for {
		select {
		case msg := <-server.Transport.Consume():
			server.Logger.Infof("%+v", msg)
		case <-server.quitChan:
			return
		}
	}
}

// Start() starts the fileserver
func (server *Fileserver) Start() error {
	if err := server.Transport.ListenAndAccept(); err != nil {
		return err
	}

	server.keepAlive()

	return nil
}

func (server *Fileserver) Stop() {
	close(server.quitChan)
	server.Logger.Info("Fileserver is closed by calling Stop() function")
	server.Logger.Close()
}

func (server *Fileserver) Store(key string, r io.Reader) error {
	return server.storage.WriteStream(key, r)
}
