package storage

type PathTransformFunc func(string) string

type StorageOptions struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}
