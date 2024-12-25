package fileserver

import (
	"io"
	"sync"

	"github.com/sumit-behera-in/goLogger"
	"github.com/sumit-behera-in/gonas/p2p"
	"github.com/sumit-behera-in/gonas/storage"
)

type FileServerOpts struct {
	Logger            *goLogger.Logger          // this is the global logger
	StorageRoot       string                    // root storage to store all files or folders managed by goNAS
	PathTransformFunc storage.PathTransformFunc // used encrypt path and file name
	Transport         p2p.Transport             // TCP, UDP, HTTP
	BootStrapNodes    []string
}

type Fileserver struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

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
		peers:          make(map[string]p2p.Peer),
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

	server.bootstrapNetwork()

	server.keepAlive()

	return nil
}

func StoreData(key string, r io.Reader) error {
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

func (server *Fileserver) bootstrapNetwork() {
	for _, address := range server.BootStrapNodes {
		go func() {
			if err := server.Transport.Dial(address); err != nil {
				server.Logger.Errorf("BootStrapping failed for address : %s", address)
			}
		}()
	}
}

func (server *Fileserver) OnPeer(p p2p.Peer) error {
	server.peerLock.Lock()
	defer server.peerLock.Unlock()
	server.peers[p.RemoteAddress().String()] = p
	server.Logger.Infof("connected with remote %s", p.RemoteAddress().String())
	return nil
}
