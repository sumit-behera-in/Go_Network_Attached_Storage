package storage

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Storage struct {
	StorageOptions
}

func NewStorage(options StorageOptions) *Storage {
	return &Storage{
		StorageOptions: options,
	}
}

// WriteStream writes the file to the disk
func (s *Storage) WriteStream(key string, r io.Reader) error {
	pathName, fileName := s.PathTransformFunc(key)

	// create the directory if doesn't exist
	if err := os.MkdirAll(pathName, os.ModeAppend); err != nil {
		return fmt.Errorf("error while creating directory %s and err: %s", pathName, err.Error())
	}
	pathWithFileName := pathName + string(filepath.Separator) + fileName

	file, err := os.Create(pathWithFileName)
	if err != nil {
		return fmt.Errorf("error while creating file %s and err: %s", fileName, err.Error())
	}

	n, err := io.Copy(file, r)
	if err != nil {
		return fmt.Errorf("error while coping file: %s", err.Error())
	}

	log.Printf("written %d bytes to disk: %s", n, pathWithFileName)

	return nil
}

// ReadStream reads the file from the disk and returns the file as a io.Reader, please close the reader after using it
func (s *Storage) ReadStream(key string) (io.Reader, error) {
	pathName, fileName := s.PathTransformFunc(key)
	pathWithFileName := pathName + string(filepath.Separator) + fileName

	file, err := os.Open(pathWithFileName)
	if err != nil {
		return nil, fmt.Errorf("error while opening file %s and err: %s", fileName, err.Error())
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	return buf, err
}
