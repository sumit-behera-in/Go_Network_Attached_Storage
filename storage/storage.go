package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Storage struct {
	StorageOptions
	mu sync.Map // Map to store mutexes for each key
}

func NewStorage(options StorageOptions) *Storage {
	return &Storage{
		StorageOptions: options,
	}
}

func (s *Storage) getMutex(key string) *sync.Mutex {
	mu, _ := s.mu.LoadOrStore(key, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

// WriteStream writes the file to the disk. The caller is responsible for closing the io.Reader passed to it.
func (s *Storage) WriteStream(key string, r io.Reader) error {
	mutex := s.getMutex(key)
	mutex.Lock()
	defer mutex.Unlock()

	pathName, fileName := s.PathTransformFunc(key)

	// create the directory if doesn't exist
	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return fmt.Errorf("error while creating directory %s and err: %s", pathName, err.Error())
	}
	pathWithFileName := pathName + string(filepath.Separator) + fileName

	file, err := os.Create(pathWithFileName)
	if err != nil {
		return fmt.Errorf("error while creating file %s and err: %s", fileName, err.Error())
	}

	_, err = io.Copy(file, r)
	if err != nil {
		return fmt.Errorf("error while copying data to file %s and err: %s", fileName, err.Error())
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("error while closing file %s and err: %s", fileName, err.Error())
	}

	return nil
}

// ReadStream reads the file from the disk and returns the file as an io.Reader.
// The caller is responsible for closing the returned io.Reader.
func (s *Storage) ReadStream(key string) (io.Reader, error) {
	mutex := s.getMutex(key)
	mutex.Lock()
	defer mutex.Unlock()

	pathName, fileName := s.PathTransformFunc(key)
	pathWithFileName := pathName + string(filepath.Separator) + fileName

	file, err := os.Open(pathWithFileName)
	if err != nil {
		return nil, fmt.Errorf("error while opening file %s and err: %s", pathWithFileName, err.Error())
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		return nil, fmt.Errorf("error while copying file content: %s", err.Error())
	}

	return buf, nil
}

func (s *Storage) Has(key string) bool {
	pathName, fileName := s.PathTransformFunc(key)
	pathWithFileName := pathName + string(filepath.Separator) + fileName

	log.Printf("Checking existence of file: %s", pathWithFileName)

	_, err := os.Stat(pathWithFileName)
	return err == nil
}

// Delete deletes the file from the disk
func (s *Storage) Delete(key string) error {
	mutex := s.getMutex(key)
	mutex.Lock()
	defer mutex.Unlock()

	if !s.Has(key) {
		return errors.New("file not found")
	}

	pathName, fileName := s.PathTransformFunc(key)
	pathWithFileName := pathName + string(filepath.Separator) + fileName

	err := os.Remove(pathWithFileName)
	if err != nil {
		return err
	}

	defer func() {
		log.Printf("deleted file: %s", pathWithFileName)
	}()

	return nil
}

func (s *Storage) CleanPath(path string) bool {
	paths := strings.Split(path, string(filepath.Separator))
	return RecursiveClean(0, "", &paths)
}
