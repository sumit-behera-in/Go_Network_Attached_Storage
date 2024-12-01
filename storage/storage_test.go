package storage

import (
	"bytes"
	"io"
	"testing"
)

var storage = Storage{
	StorageOptions: StorageOptions{
		PathTransformFunc: DefaultPathTransformFunc,
	},
}

var data = bytes.NewReader([]byte("some text"))
var key = "temp"

func TestStorage_WriteStream(t *testing.T) {
	tests := []struct {
		name    string		
		r   io.Reader
		wantErr bool
	}{
		{
			name: "write successful",
			r : data,
			wantErr: false,

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.WriteStream(key, tt.r); (err != nil) != tt.wantErr {
				t.Errorf("Storage.WriteStream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
