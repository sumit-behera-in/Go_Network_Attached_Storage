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

var WriteData = bytes.NewReader([]byte("some text"))
var WriteKey = "user1+abc.pdf"

var ReadData = bytes.NewReader([]byte("some text"))
var ReadKey = "user2+abcd.pdf"

func TestStorage_WriteStream(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "write successful",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.WriteStream(WriteKey, WriteData); (err != nil) != tt.wantErr {
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
			path, file := tt.pathTransformFunc(WriteKey)
			if path != tt.pathName {
				t.Errorf("Storage.PathTransformFunc() path does not matched wantedPath = %s, gotPath = %s", tt.pathName, path)
			}
			if file != tt.fileName {
				t.Errorf("Storage.PathTransformFunc() fileName does not matched wantedFileName = %s, gotFileName = %s", tt.fileName, file)
			}
		})
	}
}
func TestStorage_ReadStream(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "read successful",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.WriteStream(ReadKey, ReadData); err != nil {
				t.Errorf("Storage.WriteStream() error = %v", err)
			}

			// Reset ReadData before reading it again
			ReadData.Seek(0, io.SeekStart)

			reader, err := storage.ReadStream(ReadKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ReadStream() error = %v, wantErr %v", err, tt.wantErr)
			}

			r, _ := io.ReadAll(reader)
			ReadData.Seek(0, io.SeekStart)
			expectedData, _ := io.ReadAll(ReadData)

			if !bytes.Equal(r, expectedData) {
				t.Errorf("Storage.ReadStream() reader does not matched wantedReader = %v, gotReader = %v", expectedData, r)
			}
		})
	}
}
