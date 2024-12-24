package fileserver

import (
	"github.com/sumit-behera-in/goLogger"
	"github.com/sumit-behera-in/gonas/p2p"
	"github.com/sumit-behera-in/gonas/storage"
)

type FileServerOpts struct {
	Logger            *goLogger.Logger // this is the global logger 
	StorageRoot       string // root storage to store all files or folders managed by goNAS
	PathTransformFunc storage.PathTransformFunc // used encrypt path and file name
	Transport         p2p.Transport // TCP, UDP, HTTP
}

type Fileserver struct {
	FileServerOpts
}
