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
		r       io.Reader
		wantErr bool
	}{
		{
			name:    "write successful",
			r:       data,
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

func TestStorage_PathTransformFunc(t *testing.T) {
	tests := []struct {
		name              string
		key               string
		pathTransformFunc PathTransformFunc
		pathName          string
	}{
		{
			name:              "default PathTransformFunc",
			key:               "abc",
			pathTransformFunc: DefaultPathTransformFunc,
			pathName:          "abc",
		},
		{
			name:              "CAS PathTransformFunc",
			key:               "abc",
			pathTransformFunc: CASPathTransformFunc,
			pathName:          "a9993e36\\4706816a\\ba3e2571\\7850c26c\\9cd0d89d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if path := tt.pathTransformFunc(tt.key); path != tt.pathName {
				t.Errorf("Storage.PathTransformFunc() path does not matched wantedPath = %v, gotPath %v", tt.pathName, path)
			}
		})
	}
}
