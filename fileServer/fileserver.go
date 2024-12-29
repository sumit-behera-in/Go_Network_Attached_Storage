package fileserver

import (
	"bytes"
	"encoding/gob"
	"io"
	"sync"

	"github.com/sumit-behera-in/gonas/p2p"
	"github.com/sumit-behera-in/gonas/storage"
)

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
			p := Payload{}
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&p); err != nil {
				server.Logger.Fatal("Decoding failed")
			}
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
	buff := new(bytes.Buffer)
	tee := io.TeeReader(r, buff)
	if err := server.storage.WriteStream(key, tee); err != nil {
		return err
	}

	p := &Payload{
		Key:  key,
		Data: buff.Bytes(),
	}

	return server.broadcast(p)
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
	server.peers[p.RemoteAddr().String()] = p
	server.Logger.Infof("connected with remote %s", p.RemoteAddr().String())
	return nil
}

func (server *Fileserver) broadcast(p *Payload) error {
	peers := []io.Writer{} // peer contains net.conn which contains reader and writer so it is compatible with it.
	for _, peer := range server.peers {
		peers = append(peers, peer)
	}

	multiWriter := io.MultiWriter(peers...)
	return gob.NewEncoder(multiWriter).Encode(p)
}
