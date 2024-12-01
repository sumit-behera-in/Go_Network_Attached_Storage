package storage

import (
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

func (s *Storage) WriteStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	// create the directory if doesn't exist
	if err := os.MkdirAll(pathName, os.ModeAppend); err != nil {
		return fmt.Errorf("error while creating directory %s and err: %s", pathName, err.Error())
	}

	fileName := "anb.txt"
	pathWithFileName := pathName + string(filepath.Separator) + fileName

	file, err := os.Create(pathWithFileName)
	if(err != nil) {
		return fmt.Errorf("error while creating file %s and err: %s", fileName, err.Error())
	}

	n, err := io.Copy(file,r)
	if(err != nil) {
		return fmt.Errorf("error while coping file: %s",err.Error())
	}

	log.Printf("written %d bytes to disk: %s",n,pathWithFileName)

	return nil
}
