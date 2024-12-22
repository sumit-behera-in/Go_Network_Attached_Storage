package storage

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"path/filepath"
	"strings"
)

type PathTransformFunc func(string) (string, string)

type StorageOptions struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) (string, string) {
	keyContains := strings.Split(key, "+")
	userDetails := keyContains[0]
	fileName := keyContains[1]
	return userDetails, fileName
}

var CASPathTransformFunc = func(key string) (string, string) {
	keyContains := strings.Split(key, "+")
	userDetails := keyContains[0]
	fileExt := filepath.Ext(keyContains[1])
	fileName := filepath.Base(keyContains[1])

	hash := sha1.Sum([]byte(userDetails)) // convert [20] into slice use [:]
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 8
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	// encode file name
	fileNameBytes := md5.Sum([]byte(fileName))
	fileName = hex.EncodeToString(fileNameBytes[:])

	return strings.Join(paths, "/"), fileName + fileExt
}
