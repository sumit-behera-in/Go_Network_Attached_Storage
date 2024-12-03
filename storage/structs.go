package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"path/filepath"
	"strings"
)

type PathTransformFunc func(string) string

type StorageOptions struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

var CASPathTransformFunc = func(key string) string {
	hash := sha1.Sum([]byte(key)) // convert [20] into slice use [:]
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 8
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, string(filepath.Separator))
}
