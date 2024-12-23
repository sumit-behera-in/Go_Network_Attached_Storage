package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Response) error
}

type GOBDecoder struct{}

// Decode decodes the incoming data using the GOB decoder
func (dec *GOBDecoder) Decode(r io.Reader, rpc *Response) error {
	return gob.NewDecoder(r).Decode(rpc)
}

type DefaultDecoder struct{}

// Decode decodes the incoming data using the default decoder, i.e. the custom decoder
func (dec *DefaultDecoder) Decode(r io.Reader, rpc *Response) error {
	var buf bytes.Buffer // Use bytes.Buffer to accumulate the data
	temp := make([]byte, 1024)

	for {
		n, err := r.Read(temp)
		if err != nil {
			if err == io.EOF {
				// End of stream; break and return the accumulated data
				break
			}
			return err // Handle other errors
		}

		// Write the received bytes into the buffer
		buf.Write(temp[:n])

		// Telnet may send data in chunks. If you know the message ends with a specific terminator (e.g., \n), you can check here:
		if bytes.Contains(temp[:n], []byte("\n")) {
			break
		}
	}

	rpc.Payload = buf.Bytes() // Assign accumulated data to the response payload
	fmt.Println("Received data:", string(rpc.Payload))

	return nil

}
