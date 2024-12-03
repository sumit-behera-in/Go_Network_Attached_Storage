package storage

import (
	"bytes"
	"io"
	"testing"
)

var storage = Storage{
	StorageOptions: StorageOptions{
		PathTransformFunc: CASPathTransformFunc,
	},
}

var data = bytes.NewReader([]byte("some text"))
var key = "user1+abc.pdf"

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
		pathTransformFunc PathTransformFunc
		pathName          string
		fileName          string
	}{
		{
			name:              "default PathTransformFunc",
			pathTransformFunc: DefaultPathTransformFunc,
			pathName:          "user1",
			fileName:          "abc.pdf",
		},
		{
			name:              "CAS PathTransformFunc",
			pathTransformFunc: CASPathTransformFunc,
			pathName:          "b3daa77b\\4c04a955\\1b8781d0\\3191fe09\\8f325e67",
			fileName:          "c7634722815d7f16a4668d0b52f3038b.pdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, file := tt.pathTransformFunc(key)
			if path != tt.pathName {
				t.Errorf("Storage.PathTransformFunc() path does not matched wantedPath = %s, gotPath = %s", tt.pathName, path)
			}
			if file != tt.fileName {
				t.Errorf("Storage.PathTransformFunc() fileName does not matched wantedFileName = %s, gotFileName = %s", tt.fileName, file)
			}
		})
	}
}
