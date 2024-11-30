package main

import (
	"fmt"

	"github.com/sumit-behera-in/Distributed_File_Storage_GO/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")

	if err := tr.ListenAndAccept(); err != nil {
		panic(err.Error())
	}

	fmt.Println("HEllo")

	select {}
}
