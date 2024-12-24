package fileserver

func NewFileServer(option FileServerOpts) *Fileserver {
	return &Fileserver{
		FileServerOpts: option,
	}
}

// Start() starts the fileserver
func (server *Fileserver) Start() error {
	if err := server.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}
 