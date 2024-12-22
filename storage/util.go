package storage

import (
	"log"
	"os"
	"path/filepath"
)

func RecursiveClean(i int, path string, paths *[]string) bool {
	if i == len(*paths) {
		return true
	}

	path += string(filepath.Separator) + (*paths)[i]

	// If the child directory is not empty, return false
	if !RecursiveClean(i+1, path, paths) {
		return false
	}

	if ok, _ := isDirectoryEmpty(path); !ok {
		return false // directory is not empty
	}

	dir, err := os.Open(path)
	if err != nil {
		return false
	}

	// Try to read one entry from the directory
	if _, err := dir.Readdirnames(1); err == nil {
		// Directory is not empty
		log.Printf("Warn : Directory %s is not empty, not deleting \n", path)
		return false
	}

	dir.Close()

	// Directory is empty, remove it
	if err := os.RemoveAll(path); err != nil {
		return true
	}

	log.Printf("Info : Deleted directory %s \n", path)

	return true
}

// isDirectoryEmpty checks if the given directory is empty
func isDirectoryEmpty(dirPath string) (bool, error) {
	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	// Try to read one entry from the directory
	// If we can read an entry, it means the directory is not empty
	_, err = dir.Readdirnames(1) // Read one entry
	if err == nil {
		// Directory is not empty
		return false, nil
	}
	if err.Error() == "readdir: file already closed" {
		// This happens when the directory was closed during Readdirnames
		return false, nil
	}
	if err != nil {
		// Other errors (such as unexpected errors in opening or reading the directory)
		return false, err
	}

	// If no error and no entries, the directory is empty
	return true, nil
}
