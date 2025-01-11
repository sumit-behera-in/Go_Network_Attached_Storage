package storage

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"path/filepath"
	"strings"

	"github.com/sumit-behera-in/goLogger"
)

type PathTransformFunc func(string, string) (string, string)

type StorageOptions struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Logger            *goLogger.Logger
}

// DefaultPathTransformFunc is the default path transform function, which splits the key by "+"
// and returns the first part as path and the second part as file name
var DefaultPathTransformFunc = func(storageRoot string, key string) (string, string) {
	keyContains := strings.Split(key, "+")
	userDetails := keyContains[0]
	fileName := keyContains[1]
	return storageRoot + "/" + userDetails, fileName
}

// CASPathTransformFunc is a content-addressable storage (CAS) path transform function.
// It takes a key in the format "userDetails+fileName", where userDetails is a string
// representing user-specific information and fileName is the name of the file.
//
// The function performs the following transformations:
// 1. Splits the key by "+" to separate userDetails and fileName.
// 2. Computes the SHA-1 hash of the userDetails and encodes it as a hexadecimal string.
// 3. Divides the hash string into equal-sized blocks and constructs a path by joining these blocks with "/".
// 4. Computes the MD5 hash of the fileName and encodes it as a hexadecimal string.
// 5. Returns the constructed path and the encoded file name with its original extension.
//
// This function is useful for generating unique and consistent paths for storing files
// in a content-addressable storage system.
var CASPathTransformFunc = func(storageRoot string, key string) (string, string) {
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

	return storageRoot + "/" + strings.Join(paths, "/"), fileName + fileExt
}
